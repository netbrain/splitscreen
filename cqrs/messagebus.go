package cqrs

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SubscriptionType int

const (
	ViewSubscription SubscriptionType = iota
	ManagerSubscription
)

type MessageBus interface {
	Subscribe(f AggregateHandleFunc, subTyp SubscriptionType, typ ...MessageType) func()
	Emit(ctx context.Context, message ...Message) error
	Manage(ctx context.Context, message ...Message) error
}

type LocalMessageBus struct {
	mu sync.Mutex
	subscriptions map[string]map[int64][]AggregateHandleFunc
}

func NewLocalMessageBus() *LocalMessageBus {
	return &LocalMessageBus{
		subscriptions: map[string]map[int64][]AggregateHandleFunc{},
	}

}

func (l *LocalMessageBus) hashFn(s SubscriptionType, m MessageType) string {
	return fmt.Sprintf("%s-%d",m,s)
}

func (l *LocalMessageBus) Subscribe(f AggregateHandleFunc, subTyp SubscriptionType, typ ...MessageType) func() {
	l.mu.Lock()
	defer l.mu.Unlock()
	idx := time.Now().UnixNano()

	if len(typ) == 0 {
		key := l.hashFn(subTyp,"")
		if _,ok := l.subscriptions[key]; !ok {
			l.subscriptions[key] = map[int64][]AggregateHandleFunc{}
		}
		l.subscriptions[key][idx] = append(l.subscriptions[key][idx],f)
	}
	for _, t := range typ {
		key := l.hashFn(subTyp,t)
		if _,ok := l.subscriptions[key]; !ok {
			l.subscriptions[key] = map[int64][]AggregateHandleFunc{}
		}
		l.subscriptions[key][idx] = append(l.subscriptions[key][idx],f)
	}

	return func(){
		l.mu.Lock()
		defer l.mu.Unlock()

		if len(typ) == 0 {
			delete(l.subscriptions[l.hashFn(subTyp,"")],idx)
		}
		for _, t := range typ {
			delete(l.subscriptions[l.hashFn(subTyp,t)],idx)
		}
	}
}

func (l *LocalMessageBus) Manage(ctx context.Context, messages ...Message) error {
	return l.emit(ctx,ManagerSubscription,messages...)
}

func (l *LocalMessageBus) Emit(ctx context.Context, messages ...Message) error {
	return l.emit(ctx,ViewSubscription,messages...)
}

func (l *LocalMessageBus) emit(ctx context.Context, subTyp SubscriptionType, messages ...Message) error {
	for _, msg := range messages {
		l.mu.Lock()
		var subscribers []AggregateHandleFunc
		for _, s := range l.subscriptions[l.hashFn(subTyp,msg.Meta().MessageType)] {
			subscribers = append(subscribers,s...)
		}
		for _, s := range l.subscriptions[l.hashFn(subTyp,"")] {
			subscribers = append(subscribers,s...)
		}
		l.mu.Unlock()
		for _, fn := range subscribers {
			if err := fn(ctx, msg); err != nil {
				return err
			}
		}
	}
	return nil
}

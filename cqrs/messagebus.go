package cqrs

import "context"

type SubscriptionType int

const (
	ViewSubscription SubscriptionType = iota
	ManagerSubscription
)

type MessageBus interface {
	Subscribe(f AggregateHandleFunc, subTyp SubscriptionType, typ ...MessageType)
	Emit(ctx context.Context, message ...Message) error
	Manage(ctx context.Context, message ...Message) error
}

type LocalMessageBus map[SubscriptionType]map[MessageType][]AggregateHandleFunc

func NewLocalMessageBus() LocalMessageBus {
	return map[SubscriptionType]map[MessageType][]AggregateHandleFunc{
		ViewSubscription: make(map[MessageType][]AggregateHandleFunc),
		ManagerSubscription: make(map[MessageType][]AggregateHandleFunc),
	}
}

func (l LocalMessageBus) Subscribe(f AggregateHandleFunc, subTyp SubscriptionType, typ ...MessageType) {
	if len(typ) == 0 {
		l[subTyp][""] = append(l[subTyp][""], f)
	}
	for _, t := range typ {
		l[subTyp][t] = append(l[subTyp][t], f)
	}
}

func (l LocalMessageBus) Manage(ctx context.Context, messages ...Message) error {
	return l.emit(ctx,ManagerSubscription,messages...)
}

func (l LocalMessageBus) Emit(ctx context.Context, messages ...Message) error {
	return l.emit(ctx,ViewSubscription,messages...)
}

func (l LocalMessageBus) emit(ctx context.Context, subTyp SubscriptionType, messages ...Message) error {
	for _, msg := range messages {
		subscribers := l[subTyp][msg.Meta().MessageType]
		wildSubs, ok := l[subTyp][""]
		if ok {
			subscribers = append(subscribers, wildSubs...)
		}

		for _, fn := range subscribers {
			if err := fn(ctx, msg); err != nil {
				return err
			}
		}
	}
	return nil
}

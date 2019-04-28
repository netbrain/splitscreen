package cqrs

import "context"

type MessageBus interface {
	Subscribe(f AggregateHandleFunc, typ ...MessageType)
	Emit(ctx context.Context, message ...Message) error
}

type LocalMessageBus map[MessageType][]AggregateHandleFunc

func NewLocalMessageBus() LocalMessageBus {
	return make(LocalMessageBus)
}

func (l LocalMessageBus) Subscribe(f AggregateHandleFunc, typ ...MessageType) {
	if len(typ) == 0 {
		l[""] = append(l[""], f)
	}
	for _, t := range typ {
		l[t] = append(l[t], f)
	}
}

func (l LocalMessageBus) Emit(ctx context.Context, messages ...Message) error {
	for _, msg := range messages {
		subscribers := l[msg.Meta().MessageType]
		wildSubs, ok := l[""]
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

package cqrs

import (
	"context"
	"fmt"
)

const (
	TestAggregateType AggregateType = "TestAggregate"
	TestEventType     MessageType   = "TestEvent"
)

type TestAggregate struct {
	*AggregateMeta
}

func (t *TestAggregate) Register(app *App) {
	app.RegisterAggregate(TestAggregateType, func() AggregateRoot {
		return &TestAggregate{}
	})
}

func (t *TestAggregate) Meta() *AggregateMeta {
	return t.AggregateMeta
}

func (t *TestAggregate) Handle(ctx context.Context, msg Message) error {
	switch msg.Meta().MessageType {
	case TestEventType:
		return nil
	}
	return fmt.Errorf("unexpected instance")
}

type TestEvent struct {
	*MessageMeta
}

func (t *TestEvent) Meta() *MessageMeta {
	return t.MessageMeta
}

func registerTestTypes(cqrs *App) {
	cqrs.RegisterAggregate(TestAggregateType, func() AggregateRoot {
		return &TestAggregate{
			AggregateMeta: &AggregateMeta{},
		}
	})
	cqrs.RegisterMessage(func() Message {
		return &TestEvent{
			MessageMeta: &MessageMeta{
				AggregateMeta: &AggregateMeta{
					AggregateType: TestAggregateType,
				},
				MessageType: TestEventType,
			},
		}
	})
}

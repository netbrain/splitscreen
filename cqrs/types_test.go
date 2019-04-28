package cqrs

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	TestAggregateType AggregateType = "TestAggregate"
	TestEventType     MessageType   = "TestEvent"
)

func init() {
	RegisterAggregate(TestAggregateType, func() AggregateRoot {
		return &TestAggregate{
			AggregateMeta: &AggregateMeta{},
		}
	})
	RegisterMessage(func() Message {
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

type TestAggregate struct {
	*AggregateMeta
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

func (t *TestEvent) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TestEvent) Deserialize(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *TestEvent) Meta() *MessageMeta {
	return t.MessageMeta
}

func (t *TestEvent) ToRawMessage() (*RawMessage, error) {
	return NewRawMessage(t)
}

func NewTestEvent(data TestEvent, aggregateId ...string) Message {
	msg := NewMessage(TestEventType, aggregateId...)
	data.MessageMeta = msg.Meta()
	return &data
}

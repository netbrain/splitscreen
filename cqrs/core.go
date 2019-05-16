package cqrs

import (
	"context"
	"strings"
	"time"
)

type Serializer interface {
	Serialize(src interface{}) ([]byte, error)
}

type Deserializer interface {
	Deserialize(buf []byte, dst interface{}) error
}

type Factory func(typ MessageType) Message

type MessageMeta struct {
	*AggregateMeta
	MessageType MessageType `json:"type"`
	ID          string      `json:"id"`
	Timestamp   time.Time   `json:"timestamp"`
	Replay      bool
}

type MessageType string

func (m MessageType) IsCommand() bool {
	return strings.HasSuffix(string(m), "Command")
}

func (m MessageType) IsEvent() bool {
	return strings.HasSuffix(string(m), "Event")
}

type Message interface {
	Meta() *MessageMeta
}

type RawMessage struct {
	*MessageMeta
	Data []byte
}

func (e *RawMessage) Meta() *MessageMeta {
	return e.MessageMeta
}

func (e *RawMessage) ToImplementation(ctx context.Context,dst Message) error {
	app := FromContext(ctx)
	err := app.Deserialize(e.Data, dst)
	if err != nil {
		return err
	}
	m := dst.Meta()
	*m = *e.MessageMeta
	return nil
}

func NewMessage(ctx context.Context,typ MessageType, aggregateId ...string) Message {
	app := FromContext(ctx)
	msg := app.GetMessage(ctx,typ)
	if len(aggregateId) > 0 {
		msg.Meta().AggregateID = aggregateId[0]
	}
	return msg
}

func NewRawMessage(ctx context.Context,msg Message) (*RawMessage, error) {
	app := FromContext(ctx)
	buf, err := app.Serialize(msg)
	if err != nil {
		return nil, err

	}
	return &RawMessage{
		MessageMeta: msg.Meta(),
		Data:        buf,
	}, nil
}

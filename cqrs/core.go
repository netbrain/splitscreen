package cqrs

import (
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
	replay      bool
}

func (m *MessageMeta) Replayed() bool {
	return m.replay
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

func (e *RawMessage) ToImplementation(dst Message) error {
	err := Deserialize(e.Data, dst)
	if err != nil {
		return err
	}
	m := dst.Meta()
	*m = *e.MessageMeta
	return nil
}

func NewMessage(typ MessageType, aggregateId ...string) Message {
	var id string
	if len(aggregateId) == 0 || aggregateId[0] == "" {
		id = NewID()
	} else {
		id = aggregateId[0]
	}
	msg := GetMessage(typ)
	msg.Meta().AggregateID = id
	return msg
}

func NewRawMessage(msg Message) (*RawMessage, error) {
	buf, err := Serialize(msg)
	if err != nil {
		return nil, err

	}
	return &RawMessage{
		MessageMeta: msg.Meta(),
		Data:        buf,
	}, nil
}

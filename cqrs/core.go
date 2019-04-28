package cqrs

import (
	"encoding/json"
	"strings"
	"time"
)

type Serializable interface {
	Serialize() ([]byte, error)
}

type Deserializable interface {
	Deserialize([]byte) error
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
	Serializable
	Deserializable
	Meta() *MessageMeta
	ToRawMessage() (*RawMessage, error)
}

type RawMessage struct {
	*MessageMeta
	Data []byte
}

func NewRawMessage(msg Message) (*RawMessage, error) {
	buf, err := msg.Serialize()
	if err != nil {
		return nil, err
	}
	return &RawMessage{
		MessageMeta: msg.Meta(),
		Data:        buf,
	}, nil
}

func (e *RawMessage) Meta() *MessageMeta {
	return e.MessageMeta
}

func (e *RawMessage) ToRawMessage() (*RawMessage, error) {
	return e, nil
}

func (e *RawMessage) Deserialize(b []byte) error {
	return json.Unmarshal(b, e)
}

func (e *RawMessage) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *RawMessage) ToImplementation(dst Message) error {
	err := dst.Deserialize(e.Data)
	if err != nil {
		return err
	}
	m := dst.Meta()
	*m = *e.MessageMeta
	return nil
}

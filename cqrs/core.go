package cqrs

import (
	"fmt"
	"strings"
	"time"
)

var ErrMetaNotPresent = fmt.Errorf("meta not initialized on aggregate")
var ErrNoID = fmt.Errorf("no id specified on aggregate")
var ErrNoEvents = fmt.Errorf("no events")

type Serializer interface {
	Serialize(src interface{}) ([]byte, error)
}

type Deserializer interface {
	Deserialize(buf []byte, dst interface{}) error
}

type Registerable interface {
	Register(app *App)
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

func (e *RawMessage) ToImplementation(d Deserializer, dst Message) error {
	err := d.Deserialize(e.Data, dst)
	if err != nil {
		return err
	}
	m := dst.Meta()
	*m = *e.MessageMeta
	return nil
}

func NewRawMessage(s Serializer, m Message) (*RawMessage, error) {
	data, err := s.Serialize(m)
	return &RawMessage{
		MessageMeta: m.Meta(),
		Data:        data,
	}, err
}

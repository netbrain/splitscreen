package cqrs

import (
	"context"
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
	MessageType MessageType            `json:"type"`
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Replay      bool                   `json:"replay"`
	Data        map[string]interface{} `json:"data"`
}

func (m *MessageMeta) Set(key string, value interface{}) {
	if m.Data == nil {
		m.Data = map[string]interface{}{}
	}
	m.Data[key] = value
}

func (m *MessageMeta) Get(key string) interface{} {
	if m.Data == nil {
		return nil
	}
	return m.Data[key]
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

func LoadAggregate(ctx context.Context, es EventStore, meta *AggregateMeta, dst AggregateRoot) error {
	aggrMeta := dst.Meta()
	if aggrMeta == nil {
		return ErrMetaNotPresent
	}

	if aggrMeta.loaded {
		return nil
	}

	if meta.AggregateID == "" {
		return ErrNoID
	}

	result := es.Load(ctx, meta.AggregateID, meta.AggregateType)
	var count int
	for e := range result {
		count++
		if e.Err != nil {
			return e.Err
		}
		if err := dst.Handle(ctx, e.Message); err != nil {
			return err
		}
	}

	if count == 0 {
		return ErrNoEvents
	}

	*aggrMeta = *meta
	aggrMeta.loaded = true
	return nil
}

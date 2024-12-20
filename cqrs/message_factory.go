package cqrs

import (
	"reflect"
	"time"
)

type messageInfo struct {
	instance interface{}
	fn       func() Message
}

type MessageFactory interface {
	RegisterMessage(fn func() Message)
	NewMessage(typ MessageType, aggregateId ...string) Message
	NewMessageWithCause(typ MessageType, aggregateId string, causedBy *MessageMeta) Message
}

type DefaultMessageFactory struct {
	fnMap map[MessageType]*messageInfo
	idGen IDGenerator
}

func NewDefaultMessageFactory(id IDGenerator) *DefaultMessageFactory {
	return &DefaultMessageFactory{
		fnMap: make(map[MessageType]*messageInfo),
		idGen: id,
	}
}

func (a *DefaultMessageFactory) RegisterMessage(fn func() Message) {
	msg := fn()
	typ := msg.Meta().MessageType
	if _, ok := a.fnMap[typ]; ok {
		panic("aggregate type already registered!")
	}
	msgTyp := reflect.TypeOf(msg)
	if msgTyp.Kind() == reflect.Ptr {
		msgTyp = msgTyp.Elem()
	}
	v := reflect.New(msgTyp).Elem().Interface()
	a.fnMap[typ] = &messageInfo{
		instance: v,
		fn:       fn,
	}
}


func (a *DefaultMessageFactory) NewMessage(typ MessageType, aggregateId ...string) Message {
	if _, ok := a.fnMap[typ]; !ok {
		return nil
	}
	msgInfo := a.fnMap[typ]
	msg := msgInfo.fn()
	meta := msg.Meta()
	meta.ID = a.idGen.NewID()
	meta.CorrelationID = a.idGen.NewID()
	meta.Timestamp = time.Now().UTC()
	if len(aggregateId) > 0 {
		meta.AggregateID = aggregateId[0]
	}

	return msg
}

func (a *DefaultMessageFactory) NewMessageWithCause(typ MessageType, aggregateId string, causedByMeta *MessageMeta) (msg Message) {
	msg = a.NewMessage(typ,aggregateId)
	if msg == nil {
		return
	}
	meta := msg.Meta()
	meta.CausationID = causedByMeta.ID
	meta.CorrelationID = causedByMeta.CorrelationID
	return
}

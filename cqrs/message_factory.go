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
	msg.Meta().ID = a.idGen.NewID()
	msg.Meta().Timestamp = time.Now().UTC()
	if len(aggregateId) > 0 {
		msg.Meta().AggregateID = aggregateId[0]
	}

	return msg
}

package cqrs

import (
	"context"
	"reflect"
)

type AggregateFactory struct {
	fnMap map[AggregateType]func() AggregateRoot
}

func NewAggregateFactory() *AggregateFactory {
	return &AggregateFactory{
		fnMap: make(map[AggregateType]func() AggregateRoot),
	}
}

func (a *AggregateFactory) RegisterAggregate(typ AggregateType, f func() AggregateRoot) {
	if _, ok := a.fnMap[typ]; ok {
		panic("aggregate type already registered!")
	}
	a.fnMap[typ] = f
}

func (a *AggregateFactory) GetAggregate(typ AggregateType) AggregateRoot {
	if _, ok := a.fnMap[typ]; !ok {
		return nil
	}
	return a.fnMap[typ]()
}

type messageInfo struct {
	instance interface{}
	fn       func() Message
}

type MessageFactory struct {
	fnMap map[MessageType]*messageInfo
}

func NewMessageFactory() *MessageFactory {
	return &MessageFactory{
		fnMap: make(map[MessageType]*messageInfo),
	}
}

func (a *MessageFactory) RegisterMessage(fn func() Message) {
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

func (a *MessageFactory) GetMessage(ctx context.Context,typ MessageType) Message {
	app := FromContext(ctx)
	if _, ok := a.fnMap[typ]; !ok {
		return nil
	}
	msgInfo := a.fnMap[typ]
	msg := msgInfo.fn()
	msg.Meta().ID = app.NewID()
	return msg
}

package cqrs

type AggregateFactory interface {
	RegisterAggregate(typ AggregateType, f func() AggregateRoot)
	GetAggregate(typ AggregateType) AggregateRoot
}

type DefaultAggregateFactory struct {
	fnMap map[AggregateType]func() AggregateRoot
}

func NewDefaultAggregateFactory() *DefaultAggregateFactory {
	return &DefaultAggregateFactory{
		fnMap: make(map[AggregateType]func() AggregateRoot),
	}
}

func (a *DefaultAggregateFactory) RegisterAggregate(typ AggregateType, f func() AggregateRoot) {
	if _, ok := a.fnMap[typ]; ok {
		panic("aggregate type already registered!")
	}
	a.fnMap[typ] = f
}

func (a *DefaultAggregateFactory) GetAggregate(typ AggregateType) AggregateRoot {
	if _, ok := a.fnMap[typ]; !ok {
		return nil
	}
	return a.fnMap[typ]()
}

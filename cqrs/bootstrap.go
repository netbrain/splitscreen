package cqrs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"
	"reflect"
	"context"
)

type IDGenerator func() string

var (
	idFn IDGenerator = defaultIdGen
	views      = make(ViewMap)
	aggregates = make(map[AggregateType]func() *Aggregate)
	events = make(map[EventType]func() EventImpl)
	eventstore EventStore = &MemoryEventStore{}
	once       sync.Once
)

func Init(store EventStore, idFn IDGenerator) {
	once.Do(func() {
		if eventstore != nil {
			eventstore = store
		}
		if idFn != nil {
			idFn = idFn
		}
	})
}

func ValidateIDAndVersion(id string, version int) error {
	if version < 0 {
		return fmt.Errorf("invalid version %d", version)
	}

	if id == "" {
		return fmt.Errorf("no ID specified")
	}
	return nil
}

func DispatchCommand(ctx context.Context,c *Command) error {
	BroadcastCommand(c)
	if err := ValidateIDAndVersion(c.Aggregate.ID, c.Aggregate.Version); err != nil {
		return err
	}
	a, err := CreateOrLoadAggregate(ctx,&c.Aggregate)
	if err != nil {
		return err
	}

	if err := a.Handle(ctx,c); err != nil {
		return err
	}
	return a.Commit()
}

func CreateOrLoadAggregate(ctx context.Context,base *AggregateID) (*Aggregate, error) {
	aggregate := aggregates[base.Type]()
	aggregate.replayMode = true
	for _, e := range eventstore.Load(base.ID, base.Type) {
		if err := aggregate.Apply(ctx,e); err != nil {
			return nil, err
		}
	}
	aggregate.replayMode = false
	return aggregate, nil
}

func RegisterAggregate(typ AggregateType, newFunc func() *Aggregate) {
	aggregates[typ] = newFunc
}

func RegisterView(typ ViewType, v interface{}) {
	views[typ] = v
}

func RegisterEvent(typ EventType, v EventImpl){
	if reflect.ValueOf(v).Kind() != reflect.Struct{
		panic("expected a struct")
	}
	events[typ] = func() EventImpl{
		x := &v
		v = *x
		return v
	}

}

func View(typ ViewType) interface{} {
	v, ok := views[typ]
	if !ok {
		panic("no view of type: " + typ + " registered")
	}
	return v
}

func defaultIdGen() string {
	str := fmt.Sprintf("%s%d", time.Now().Format(time.StampNano), rand.Int())
	buf := []byte(str)
	sum256 := sha256.Sum256(buf)
	return hex.EncodeToString(sum256[:])

}

func IDFunc() string {
	return idFn()
}

package cqrs

import (
	"time"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
)

var(
	IDFunc =  func() string {
		rand.Seed(time.Now().UnixNano())
		str := fmt.Sprintf("%s%d",time.Now().Format(time.StampNano),rand.Int())
		buf := []byte(str)
		sum256 := sha256.Sum256(buf)
		return hex.EncodeToString(sum256[:])
	}
	views      = make(ViewMap)
 	aggregates = make(map[AggregateType]func()*Aggregate)
)


func ValidateIDAndVersion(id string, version int) error {
	if version < 0 {
		return fmt.Errorf("invalid version %d",version)
	}

	if id == "" {
		return fmt.Errorf("no ID specified")
	}
	return nil
}

func DispatchCommand(c *Command) error {
	BroadcastCommand(c)
	if err := ValidateIDAndVersion(c.Aggregate.ID,c.Aggregate.Version); err != nil {
		return err
	}
	a,err := CreateOrLoadAggregate(&c.Aggregate)
	if err != nil {
		return err
	}

	if err := a.Handle(c); err != nil {
		return err
	}
	return a.Commit()
}

func CreateOrLoadAggregate(base *AggregateID) (*Aggregate,error) {
	aggregate := aggregates[base.Type]()
	aggregate.replayMode = true
	for _, e := range Load(base.ID, base.Type) {
		if err := aggregate.Apply(e); err != nil {
			return nil,err
		}
	}
	aggregate.replayMode = false
	return aggregate,nil
}

func RegisterAggregate(typ AggregateType, newFunc func()*Aggregate){
	aggregates[typ] = newFunc
}

func RegisterView(typ ViewType, v interface{}){
	views[typ] = v
}

func View(typ ViewType) interface{} {
	v,ok := views[typ]
	if !ok {
		panic("no view of type: "+typ+" registered")
	}
	return v
}
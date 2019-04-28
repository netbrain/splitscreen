package cqrs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

type IDGenerator interface {
	NewID() string
}

type DefaultIDGenerator struct {
	rand *rand.Rand
}

func (d *DefaultIDGenerator) NewID() string {
	str := fmt.Sprintf("%s%d", time.Now().Format(time.StampNano), d.rand.Int())
	buf := []byte(str)
	sum256 := sha256.Sum256(buf)
	return hex.EncodeToString(sum256[:])
}

func NewDefaultIDGenerator() *DefaultIDGenerator {
	return &DefaultIDGenerator{rand: rand.New(rand.NewSource(time.Now().Unix()))}
}

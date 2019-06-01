package json

import "encoding/json"

type Serializer struct{}

func (*Serializer) Serialize(src interface{}) ([]byte, error) {
	return json.Marshal(src)
}

type Deserializer struct{}

func (*Deserializer) Deserialize(buf []byte, dst interface{}) error {
	return json.Unmarshal(buf, dst)
}

func New() (*Serializer, *Deserializer) {
	return &Serializer{}, &Deserializer{}
}

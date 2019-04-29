package json

import "encoding/json"

type Serializer struct{}

func NewSerializer() *Serializer {
	return &Serializer{}
}

func (*Serializer) Serialize(src interface{}) ([]byte, error) {
	return json.Marshal(src)
}

type Deserializer struct{}

func NewDeserializer() *Deserializer {
	return &Deserializer{}
}

func (*Deserializer) Deserialize(buf []byte, dst interface{}) error {
	return json.Unmarshal(buf, dst)
}

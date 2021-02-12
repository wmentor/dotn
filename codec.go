package dotn

import (
	"encoding/json"
)

// used for create object from complex structs
type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type jsonCodec struct {
}

func (j *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func NewJsonCodec() Codec {
	return &jsonCodec{}
}

package config

import (
	"encoding/json"
)

type Output struct {
	Type  string `bson:"resourceType"`
	Name  string `bson:"resourceName"`
	Value interface{}
}

type OutputList []Output

type OutputPayload struct {
	Name  string `bson:"name"`
	Type  string `bson:"resourceType"`
	Value string `bson:"value"`
}

type OutputPayloadList []OutputPayload

func (o *Output) ToString() string {
	b, _ := json.Marshal(o)
	return string(b)
}

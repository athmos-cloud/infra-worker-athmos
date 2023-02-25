package config

import "reflect"

type VariableType string

type Input struct {
	Name        string       `bson:"name"`
	Description string       `bson:"description"`
	Type        reflect.Kind `bson:"type"`
	Default     string       `bson:"default"`
	Value       interface{}  `bson:"value"`
}

type InputList []Input

type InputPayload struct {
	Name  string      `bson:"name"`
	Value interface{} `bson:"value"`
}

type InputPayloadList []InputPayload

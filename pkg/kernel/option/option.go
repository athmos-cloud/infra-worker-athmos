package option

import (
	"reflect"
)

type Option struct {
	Type  string
	Value interface{}
}

func New(kind string, v ...interface{}) Option {
	if kind != reflect.TypeOf(reflect.TypeOf(v).Kind()).String() {
		return Option{
			Type:  reflect.Zero(reflect.TypeOf(v)).String(),
			Value: v,
		}
	}
	if len(v) == 0 {
		return Option{
			Type:  kind,
			Value: []interface{}{},
		}
	}
	return Option{
		Type:  kind,
		Value: v,
	}
}

func Null() Option {
	return Option{}
}

func EmptyWithMessage(msg string) *Option {
	return &Option{
		Value: msg,
	}
}

func (o Option) Validate() bool {
	if o.Value == nil {
		return false
	}
	return o.Type == reflect.TypeOf(o.Value).String()
}

func (o Option) SetType(t string) Option {
	o.Type = t
	return o
}

func (o Option) Get() interface{} {
	return o.Value
}

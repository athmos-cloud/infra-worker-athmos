package option

import "reflect"

type Option struct {
	Type  reflect.Kind
	Value interface{}
}

func New(kind reflect.Kind, v ...interface{}) Option {
	if kind != reflect.TypeOf(reflect.TypeOf(v).Kind()).Kind() {
		return Option{
			Type:  reflect.Zero(reflect.TypeOf(v)).Kind(),
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
	return reflect.DeepEqual(o.Type, reflect.TypeOf(o.Value).Kind())
}

func (o Option) SetType(t reflect.Kind) Option {
	o.Type = t
	return o
}

func (o Option) Get() interface{} {
	return o.Value
}

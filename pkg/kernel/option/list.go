package option

import (
	"reflect"
)

type List struct {
	Type   reflect.Kind
	Values []interface{}
}

func NewList(kind reflect.Kind, args ...interface{}) *List {
	var values []interface{}
	for _, arg := range args {
		values = append(values, arg)
	}
	return &List{
		Type:   kind,
		Values: values,
	}
}

func (ol *List) Validate(size ...int) bool {
	if len(size) > 0 && len(ol.Values) != size[0] {
		return false
	}
	for _, o := range ol.Values {
		if reflect.TypeOf(o).Kind() != ol.Type {
			return false
		}
	}
	return true
}

func (ol *List) SetType(t reflect.Kind) *List {
	ol.Type = t
	return ol
}

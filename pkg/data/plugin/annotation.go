package plugin

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/auth"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
	"strings"
)

const (
	tagName         = "plugin"
	nestedSeparator = "."
	authTypeKey     = "authType"
	diskModeKey     = "diskMode"
)

func InjectMapIntoStruct(m map[string]interface{}, s interface{}) errors.Error {
	structValue := reflect.ValueOf(s).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		tag := structType.Field(i).Tag.Get(tagName)
		if tag == "" {
			continue
		}
		keys := strings.Split(tag, nestedSeparator)
		value, ok := getNestedValue(m, keys)
		if !ok {
			continue
		}
		fieldValue := reflect.ValueOf(value)
		fieldType := field.Type()
		if fieldType.Kind() == reflect.Struct {
			nestedStructPtr := field.Addr().Interface()
			err := InjectMapIntoStruct(value.(map[string]interface{}), nestedStructPtr)
			if !err.IsOk() {
				return err
			}
		} else if fieldType == fieldValue.Type() {
			field.Set(fieldValue)
		} else if fieldType.Kind() == reflect.Slice && fieldValue.Type().Kind() == reflect.Slice {
			// Handle slice of structs
			slice := reflect.MakeSlice(fieldType, fieldValue.Len(), fieldValue.Len())
			for j := 0; j < fieldValue.Len(); j++ {
				if fieldValue.Index(j).Type().Kind() == reflect.Map {
					nestedMap := fieldValue.Index(j).Interface().(map[string]interface{})
					nestedStruct := reflect.New(fieldType.Elem()) // Create a new instance of the struct
					err := InjectMapIntoStruct(nestedMap, nestedStruct.Interface())
					if !err.IsOk() {
						return err
					}
					slice.Index(j).Set(nestedStruct.Elem()) // Set the struct value in the slice
				} else {
					convertedValue := reflect.ValueOf(fieldValue.Index(j).Interface()).Convert(fieldType.Elem())
					slice.Index(j).Set(convertedValue)
				}
			}
			field.Set(slice)
		} else {
			if err := handleEnumTypes(keys[0], field, fieldValue); !err.IsOk() {
				return err
			}
		}
	}

	return errors.OK
}

func handleEnumTypes(key string, field reflect.Value, value reflect.Value) errors.Error {
	if key == authTypeKey {
		if _, err := auth.AuthType(value.String()); !err.IsOk() {
			return err
		}
		field.SetString(value.String())
	} else if key == diskModeKey {
		if res := types.DiskModeType(value.String()); res == "" {
			return errors.InvalidArgument.WithMessage(fmt.Sprintf("invalid disk mode: %s", value.String()))
		}
		field.SetString(value.String())
	} else {
		return errors.InvalidArgument.WithMessage(fmt.Sprintf("type mismatch for field: %s", key))
	}
	return errors.OK
}

func getNestedValue(m map[string]interface{}, keys []string) (interface{}, bool) {
	value := m[keys[0]]

	if len(keys) == 1 {
		return value, value != nil
	}

	if nestedMap, ok := value.(map[string]interface{}); ok {
		return getNestedValue(nestedMap, keys[1:])
	}

	return nil, false
}

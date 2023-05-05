package mongo

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

func parseMongoEntry(entry interface{}) bson.M {
	bytes, err := bson.Marshal(entry)
	if err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
	var result map[string]interface{}
	if err = bson.Unmarshal(bytes, &result); err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
	return parseBsonMap(result)
}

func parseBsonMap(aMap map[string]interface{}) bson.M {
	finalMap := make(map[string]interface{})
	parseMap("", aMap, &finalMap)
	return finalMap
}

func parseMap(k string, aMap map[string]interface{}, finalMap *map[string]interface{}) {
	if len(aMap) == 0 {
		(*finalMap)[k] = nil
		return
	}

	for key, val := range aMap {
		if val != nil {
			switch concreteVal := val.(type) {
			case map[string]interface{}:
				parseMap(getKey(k, key), val.(map[string]interface{}), finalMap)
			case []interface{}:
				(*finalMap)[getKey(k, key)] = val.([]interface{})
			default:
				concreteValType := reflect.TypeOf(concreteVal)
				if concreteValType.Kind() == reflect.Map {
					parseMap(getKey(k, key), concreteVal.(primitive.M), finalMap)
				} else {
					(*finalMap)[getKey(k, key)] = concreteVal
				}
			}
		} else {
			(*finalMap)[getKey(k, key)] = nil
		}
	}
}

func getKey(k string, key string) string {
	if k == "" {
		return key
	}
	return k + "." + key
}

func generateUpdateBSON(prefix string, value interface{}) bson.M {
	update := bson.M{}
	v := reflect.ValueOf(value)

	processStruct := func(structValue reflect.Value, structType reflect.Type) {
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			fieldValue := structValue.Field(i)

			bsonTag := field.Tag.Get("bson")
			bsonFieldName := strings.Split(bsonTag, ",")[0]

			if bsonFieldName == "" || bsonFieldName == "-" {
				continue
			}

			fullName := prefix + bsonFieldName

			if fieldValue.Kind() == reflect.Struct {
				for k, v := range generateUpdateBSON(fullName+".", fieldValue.Interface()) {
					update[k] = v
				}
			} else if fieldValue.Kind() == reflect.Map {
				for k, v := range generateUpdateBSON(fullName+".", fieldValue.Interface()) {
					update[k] = v
				}
			} else {
				update[fullName] = fieldValue.Interface()
			}
		}
	}

	processMap := func(mapValue reflect.Value) {
		for _, key := range mapValue.MapKeys() {
			fieldValue := mapValue.MapIndex(key)
			fullName := prefix + key.String()

			if fieldValue.Kind() == reflect.Struct {
				for k, v := range generateUpdateBSON(fullName+".", fieldValue.Interface()) {
					update[k] = v
				}
			} else if fieldValue.Kind() == reflect.Map {
				for k, v := range generateUpdateBSON(fullName+".", fieldValue.Interface()) {
					update[k] = v
				}
			} else {
				update[fullName] = fieldValue.Interface()
			}
		}
	}

	if v.Kind() == reflect.Struct {
		processStruct(v, v.Type())
	} else if v.Kind() == reflect.Map {
		processMap(v)
	}

	return update
}

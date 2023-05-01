package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

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

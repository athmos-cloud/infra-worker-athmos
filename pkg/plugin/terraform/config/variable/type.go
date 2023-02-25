package variable

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"reflect"
	"strconv"
)

type TerraformVariableType string

func TypeFromString(value string) reflect.Kind {
	switch value {
	case "string":
		return reflect.String
	case "number":
		return reflect.Int
	case "bool":
		return reflect.Bool
	case "list":
		return reflect.TypeOf([]interface{}{}).Kind()
	case "map":
		return reflect.TypeOf(map[string]interface{}{}).Kind()
	default:
		return reflect.Zero(reflect.TypeOf(value)).Kind()
	}
}

func ToTerraformVariableValue(value interface{}) string {
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return fmt.Sprintf("\"%s\"", value.(string))
	case reflect.Int:
		return fmt.Sprintf("%d", value.(int))
	case reflect.Float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(value.(bool))
	case reflect.TypeOf([]interface{}{}).Kind():
		res := "[\n"
		for _, v := range value.([]interface{}) {
			res += fmt.Sprintf("\t%s,\n", ToTerraformVariableValue(v))
		}
		return res + "\n]"
	case reflect.TypeOf(map[string]interface{}{}).Kind():
		v := value.(map[string]interface{})
		res := "{\n"
		for key, val := range v {
			res += fmt.Sprintf("\t%s = %s\n", key, ToTerraformVariableValue(val))
		}
		return res + "\n}"
	default:
		logger.Warning.Printf("Unhandled terraform var type %s", reflect.TypeOf(value).Kind())
		return "{}"
	}
}

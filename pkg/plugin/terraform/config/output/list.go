package output

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"reflect"
)

type List []Output

func (ol *List) Merge(output Output, optn option.Option) *List {
	if !optn.SetType(reflect.Bool).Validate() {
		logger.Warning.Printf(
			fmt.Sprintf("Invalid option type. Expected Outputs, got :  %s", reflect.TypeOf(optn.Value).Kind()),
		)
	}
	for i, op := range *ol {
		if op.Hash() == output.Hash() {
			if optn.Value.(bool) {
				(*ol)[i] = output
			}
			return ol
		}
	}
	return ol
}

func (ol *List) MergeList(outputs List, optn option.Option) *List {
	if !optn.SetType(reflect.Bool).Validate() {
		logger.Warning.Printf(
			fmt.Sprintf("Invalid option type. Expected Outputs, got :  %s", reflect.TypeOf(optn.Value).Kind()),
		)
	}

	for _, output := range outputs {
		ol.Merge(output, optn)
	}
	return ol
}

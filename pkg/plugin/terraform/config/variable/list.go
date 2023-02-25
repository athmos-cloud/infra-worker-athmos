package variable

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"
	"reflect"
)

type List []Variable

func ListFromCommon(commonList config.InputList) *List {
	list := List{}
	for _, commonVariable := range commonList {
		list = append(list, *FromCommon(commonVariable))
	}
	return &list
}

func (vl *List) Merge(variable Variable, optn option.Option) *List {
	if !optn.SetType(reflect.Bool).Validate() {
		logger.Warning.Printf(
			fmt.Sprintf("Invalid option type. Expected Outputs, got :  %s", reflect.TypeOf(optn.Value).Kind()),
		)
	}
	for i, op := range *vl {
		if op.Hash() == variable.Hash() {
			if optn.Value.(bool) {
				(*vl)[i] = variable
			}
			return vl
		}
	}
	return vl
}

func (vl *List) MergeList(variables List, optn option.Option) *List {
	if !optn.SetType(reflect.Bool).Validate() {
		logger.Warning.Printf(
			fmt.Sprintf("Invalid option type. Expected Outputs, got :  %s", reflect.TypeOf(optn.Value).Kind()),
		)
	}
	for _, variable := range variables {
		vl.Merge(variable, optn)
	}
	return vl
}

func (vl *List) ToString() string {
	res := ""
	for _, variable := range *vl {
		res += fmt.Sprintf("%s\n", variable.ToString())
	}
	return res
}

func (vl *List) Build(claimName string, inputs config.InputPayloadList) *PackagedList {
	packagedList := PackagedList{}
	for _, variable := range *vl {
		for _, input := range inputs {
			if input.Name == variable.Name {
				packagedList = append(packagedList, *variable.Build(claimName, input))
			}
		}
	}
	return &packagedList
}

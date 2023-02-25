package module

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config/variable"
	"reflect"
)

type Module struct {
	Name      string               `hcl:"name,label"`
	Version   string               `hcl:"version,optional"`
	Source    string               `hcl:"source,attr"`
	Variables []variable.Reference `hcl:"values,block"`
}

type ModuleList []Module

func GetModuleName(pluginName string, callName string) string {
	return fmt.Sprintf("%s-%s", pluginName, callName)
}

func (m *Module) Hash() string {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		logger.Warning.Printf("Error encoding output %s", err)
		return ""
	}
	return string(b.Bytes())
}

func (m *Module) ToString() string {
	inputVars := ""
	for _, v := range m.Variables {
		inputVars += fmt.Sprintf("%s\n", v.ToString())
	}
	return `
		module "` + m.Name + `" {
			source = "` + m.Source + `"
			` + inputVars + `
		}
	`
}

func (ml *ModuleList) Merge(module Module, optn option.Option) *ModuleList {
	if !optn.SetType(reflect.Bool).Validate() {
		logger.Warning.Printf(
			fmt.Sprintf("Invalid option type. Expected Outputs, got :  %s", reflect.TypeOf(optn.Value).Kind()),
		)
	}
	for i, op := range *ml {
		if op.Hash() == module.Hash() {
			if optn.Value.(bool) {
				(*ml)[i] = module
			}
			return ml
		}
	}
	return ml
}

func (ml *ModuleList) MergeList(modules ModuleList, optn option.Option) *ModuleList {
	if !optn.SetType(reflect.Bool).Validate() {
		logger.Warning.Printf(
			fmt.Sprintf("Invalid option type. Expected Outputs, got :  %s", reflect.TypeOf(optn.Value).Kind()),
		)
	}
	for _, module := range modules {
		ml.Merge(module, optn)
	}
	return ml
}

func (ml *ModuleList) ToString() string {
	output := ""
	for _, m := range *ml {
		output += fmt.Sprintf("%s\n", m.ToString())
	}
	return output
}

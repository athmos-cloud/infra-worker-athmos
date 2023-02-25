package variable

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	common "github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"reflect"
)

type Variable struct {
	Name        string         `hcl:",label"`
	Description string         `hcl:"description,optional"`
	Sensitive   bool           `hcl:"sensitive,optional"`
	Type        *hcl.Attribute `hcl:"type,optional"`
	Default     *hcl.Attribute `hcl:"default,optional"`
	Options     hcl.Body       `hcl:",remain"`
	Value       interface{}
}

type Validation struct {
	Condition    string `hcl:"condition,optional"`
	ErrorMessage string `hcl:"error_message,optional"`
}

func FromCommon(input common.Input) *Variable {
	exprDefault, err := hclsyntax.ParseExpression([]byte(input.Default), "_.hcl", hcl.Pos{Line: 1, Column: 1})
	if err != nil {
		logger.Warning.Printf("Error parsing expression %s", err)
	}
	exprType, err := hclsyntax.ParseExpression([]byte(input.Type.String()), "_.hcl", hcl.Pos{Line: 1, Column: 1})
	if err != nil {
		logger.Warning.Printf("Error parsing expression %s", err)
	}
	return &Variable{
		Name:        input.Name,
		Description: input.Description,
		Default: &hcl.Attribute{
			Name: "default",
			Expr: exprDefault,
		},
		Type: &hcl.Attribute{
			Name: "type",
			Expr: exprType,
		},
	}
}

func (v *Variable) Hash() string {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(v)
	if err != nil {
		logger.Warning.Printf("Error encoding output %s", err)
		return ""
	}
	return string(b.Bytes())
}

func (v *Variable) ToReference(source string) *Reference {
	return &Reference{
		Name:   v.Name,
		Source: source,
	}
}

func (v *Variable) ToCommon() *common.Input {
	return &common.Input{
		Name:        v.Name,
		Description: v.Description,
		Default:     ToTerraformVariableValue(v.Default),
		Value:       v.Value,
		Type:        reflect.TypeOf(v.Value).Kind(),
	}
}

func (v *Variable) Build(claimName string, input common.InputPayload) *Packaged {
	variableName := fmt.Sprintf("%s_%s", claimName, v.Name)
	newV := v
	newV.Name = variableName

	return &Packaged{
		Variable:        newV,
		ModuleReference: newV.ToReference(fmt.Sprintf("var.%s", variableName)),
		TFVar: &Reference{
			Name:   variableName,
			Source: ToTerraformVariableValue(input.Value),
		},
	}
}

func (v *Variable) ToString() string {
	var value string
	if v.Value == nil {
		value = ToTerraformVariableValue(v.Default)
	} else {
		value = ToTerraformVariableValue(v.Value)
	}
	return fmt.Sprintf("%s = %s", v.Name, value)
}

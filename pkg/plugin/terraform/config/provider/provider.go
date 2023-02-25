package provider

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	common "github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"
	"reflect"
)

type Provider struct {
	Name    string `hcl:"name,label"`
	Version string `hcl:"version,optional"`
}

type List []Provider

func FromCommon(provider common.Dependency) *Provider {
	return &Provider{
		Name:    provider.Name,
		Version: provider.Version,
	}
}

func (p *Provider) Hash() string {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(p)
	if err != nil {
		logger.Warning.Printf("Error encoding output %s", err)
		return ""
	}
	return string(b.Bytes())
}

func (p *Provider) ToCommon() *common.DependencyCall {
	return &common.DependencyCall{
		Name:    p.Name,
		Version: p.Version,
	}
}

func (p *Provider) ToString() string {
	return `
		provider "` + p.Name + `" {
			version = "` + p.Version + `"
		}
	`
}

func (pl *List) Merge(provider Provider, optn option.Option) *List {
	if !optn.SetType(reflect.Bool).Validate() {
		logger.Warning.Printf(
			fmt.Sprintf("Invalid option type. Expected Outputs, got :  %s", reflect.TypeOf(optn.Value).Kind()),
		)
	}
	for i, op := range *pl {
		if op.Hash() == provider.Hash() {
			if optn.Value.(bool) {
				(*pl)[i] = provider
			}
			return pl
		}
	}
	return pl
}

func (pl *List) MergeList(providers List, optn option.Option) *List {
	if !optn.SetType(reflect.Bool).Validate() {
		logger.Warning.Printf(
			fmt.Sprintf("Invalid option type. Expected Outputs, got :  %s", reflect.TypeOf(optn.Value).Kind()),
		)
	}
	for _, module := range providers {
		pl.Merge(module, optn)
	}
	return pl
}

func (pl *List) ToString() string {
	output := ""
	for _, p := range *pl {
		output += fmt.Sprintf("%s\n", p.ToString())
	}
	return output
}

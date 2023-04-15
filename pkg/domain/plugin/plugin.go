package plugin

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/common"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

const (
	MainPluginFile = "main.yaml"
	TypePluginFile = "types.yaml"
)

type Plugin struct {
	Prerequisites []Prerequisite `yaml:"prerequisites"`
	Inputs        []Input        `yaml:"inputs"`
	Types         []Type         `yaml:"metadata,omitempty"`
}

type Prerequisite struct {
	Message   string    `yaml:"message"`
	Action    string    `yaml:"action"`
	Condition Condition `yaml:"condition"`
	Values    []string  `yaml:"with_values"`
}

type Condition struct {
	Assert string      `yaml:"assert"`
	Equals interface{} `yaml:"equals"`
}

type Input struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description,omitempty"`
	Type        string      `yaml:"type" default:"string"`
	Default     interface{} `yaml:"default,omitempty"`
	Required    bool        `yaml:"required,omitempty" default:"false"`
}

type Type struct {
	Name   string           `yaml:"name"`
	Fields map[string]Input `yaml:"fields"`
}

func Get(provider common.ProviderType, resourceType common.ResourceType) (Plugin, errors.Error) {
	//read plugin
	mainPath := fmt.Sprintf("%s/%s/%s/%s", config.Current.Plugins.Location, provider, resourceType, MainPluginFile)
	pluginBytes, err := os.ReadFile(mainPath)
	if err != nil {
		return Plugin{}, errors.IOError.WithMessage(err.Error())
	}
	plugin := Plugin{}
	if err = yaml.Unmarshal(pluginBytes, &plugin); err != nil {
		return Plugin{}, errors.ConversionError.WithMessage(err.Error())
	}
	typePath := fmt.Sprintf("%s/%s/%s/%s", config.Current.Plugins.Location, provider, resourceType, TypePluginFile)
	if _, errExists := os.Stat(typePath); os.IsNotExist(errExists) {
		return plugin, errors.OK
	}
	typesBytes, err := os.ReadFile(typePath)
	if err != nil {
		return Plugin{}, errors.IOError.WithMessage(err.Error())
	}
	if err = yaml.Unmarshal(typesBytes, &plugin.Types); err != nil {
		return Plugin{}, errors.ConversionError.WithMessage(err.Error())
	}
	return plugin, errors.OK
}

func validateMetadataPlugin(entry map[string]interface{}) (map[string]interface{}, errors.Error) {
	if entry["name"] == nil {
		return entry, errors.InvalidArgument.WithMessage("Expected name to be set")
	}
	if entry["monitored"] == nil || reflect.TypeOf(entry["monitored"]).Kind() != reflect.Bool {
		entry["monitored"] = true
	}
	if entry["tags"] == nil || reflect.TypeOf(entry["tags"]).Kind() != reflect.Map {
		entry["tags"] = map[string]string{}
	}
	return entry, errors.OK
}

func (p *Plugin) ValidateAndComplete(entry map[string]interface{}) (map[string]interface{}, errors.Error) {
	entry, err := validateMetadataPlugin(entry)
	if !err.IsOk() {
		return entry, err
	}
	for _, input := range p.Inputs {
		if entry[input.Name] == nil && input.Required && input.Default == nil {
			return entry, errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be set", input.Name))
		}
		if err2 := input.Validate(entry, p.Types); !err2.IsOk() {
			return entry, err2
		}
	}
	return entry, errors.OK
}

func (i Input) Validate(entry map[string]interface{}, types []Type) errors.Error {
	validateVar := func(input Input, entry map[string]interface{}) errors.Error {
		val := entry[input.Name]
		if val == nil && input.Default == nil && input.Required {
			return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be set", input.Name))
		} else if val == nil && (input.Default != nil || !input.Required) {
			return errors.OK
		}
		switch input.Type {
		case "string":
			if reflect.TypeOf(val).Kind() != reflect.String {
				return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a string", input.Name))
			}
			return errors.OK
		case "int":
			if reflect.TypeOf(val).Kind() != reflect.Int {
				return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be an int", input.Name))
			}
			return errors.OK
		case "bool":
			if reflect.TypeOf(val).Kind() != reflect.Bool {
				return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a bool", input.Name))
			}
			return errors.OK
		case "float":
			if reflect.TypeOf(val).Kind() != reflect.Float64 {
				return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a float", input.Name))
			}
			return errors.OK
		case "list":
			if reflect.TypeOf(val).Kind() != reflect.Slice {
				return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a list", input.Name))
			}
			return errors.OK
		}
		return errors.ValidationError.WithMessage(fmt.Sprintf("%s is not a primary type", input.Type))
	}
	if err := validateVar(i, entry); err.IsOk() {
		return errors.OK
	}
	for _, t := range types {
		if t.Name == i.Type {
			for name, input := range t.Fields {
				input.Name = name
				if err := validateVar(input, entry[t.Name].(map[string]interface{})); !err.IsOk() {
					return err
				}
			}
			return errors.OK
		} else if i.Type == fmt.Sprintf("list[%s]", t.Name) {
			subEntry := entry[i.Name].([]map[string]interface{})
			for _, sub := range subEntry {
				for _, input := range t.Fields {
					if err := validateVar(input, sub); !err.IsOk() {
						return err
					}
				}
			}
			return errors.OK
		}
	}
	return errors.OK
}

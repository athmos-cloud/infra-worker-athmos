package plugin

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/common"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

const (
	MainPluginFile = "main.yaml"
)

type Plugin struct {
	Prerequisites []Prerequisite `yaml:"prerequisites"`
	Inputs        []Input        `yaml:"inputs"`
	Types         []Type         `yaml:"types,omitempty"`
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
	Description string      `yaml:"description"`
	Type        string      `yaml:"type" default:"string"`
	Default     interface{} `yaml:"default,omitempty"`
	Required    bool        `yaml:"required" default:"false"`
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
	typePath := fmt.Sprintf("%s/%s/%s/types.yaml", config.Current.Plugins.Location, provider, resourceType)
	typesBytes, err := os.ReadFile(typePath)
	if err != nil {
		return Plugin{}, errors.IOError.WithMessage(err.Error())
	}
	if err = yaml.Unmarshal(typesBytes, &plugin.Types); err != nil {
		return Plugin{}, errors.ConversionError.WithMessage(err.Error())
	}
	return plugin, errors.OK
}

func (p *Plugin) Validate(entry map[string]interface{}) errors.Error {
	for _, input := range p.Inputs {
		if err := input.Validate(entry, p.Types); !err.IsOk() {
			return err
		}
	}
	return errors.OK
}

func (i *Input) Validate(entry map[string]interface{}, types []Type) errors.Error {
	validateType := func(subEntry map[string]interface{}) errors.Error {
		for key, val := range subEntry {
			if key == i.Name {
				switch i.Type {
				case "string":
					if reflect.TypeOf(val).Kind() != reflect.String {
						return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a string", key))
					}
				case "int":
					if reflect.TypeOf(val).Kind() != reflect.Int {
						return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be an int", key))
					}
				case "bool":
					if reflect.TypeOf(val).Kind() != reflect.Bool {
						return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a bool", key))
					}
				case "float":
					if reflect.TypeOf(val).Kind() != reflect.Float64 {
						return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a float", key))
					}
				case "list":
					if reflect.TypeOf(val).Kind() != reflect.Slice {
						return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a list", key))
					}
				default:
					return errors.ValidationError.WithMessage(fmt.Sprintf("Unknown type %s", i.Type))
				}
			}
		}
		return errors.OK
	}
	for key, val := range entry {
		if key == i.Name {
			switch i.Type {
			case "string":
				if reflect.TypeOf(val).Kind() != reflect.String {
					return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a string", key))
				}
			case "int":
				if reflect.TypeOf(val).Kind() != reflect.Int {
					return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be an int", key))
				}
			case "bool":
				if reflect.TypeOf(val).Kind() != reflect.Bool {
					return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a bool", key))
				}
			case "float":
				if reflect.TypeOf(val).Kind() != reflect.Float64 {
					return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a float", key))
				}
			case "list":
				if reflect.TypeOf(val).Kind() != reflect.Slice {
					return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a list", key))
				}
			default:
				// Check custom types
				for _, t := range types {
					if i.Type == t.Name {
						subentry := val.(map[string]interface{})
						for name, input := range t.Fields {
							for subKey, subVal := range subentry {
								if subKey == name {
									if reflect.TypeOf(subVal).String() != input.Type || subentry[name] == nil {
										return errors.ValidationError.WithMessage(fmt.Sprintf("Expected %s to be a %s", subKey, input.Type))
									}
									logger.Info.Println(subentry)
								}
							}
						}
					} else if i.Type == fmt.Sprintf("list[%s]", t.Name) {

					} else {
						return errors.ValidationError.WithMessage(fmt.Sprintf("Unknown type %s", i.Type))
					}
				}
			}

		}
	}
	return errors.OK
}

package plugin

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/kamva/mgm/v3"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

const (
	MainPluginFile = "main.yaml"
	TypePluginFile = "types.yaml"
)

type Plugin struct {
	mgm.DefaultModel `bson:",inline"`
	Prerequisites    []Prerequisite `yaml:"prerequisites" bson:"prerequisites"`
	Inputs           []Input        `yaml:"inputs" bson:"inputs"`
	Types            []Type         `yaml:"metadata,omitempty" bson:"types"`
}

type Prerequisite struct {
	mgm.DefaultModel `bson:",inline"`
	Message          string    `yaml:"message" bson:"message"`
	Action           string    `yaml:"action" bson:"action"`
	Condition        Condition `yaml:"condition" bson:"condition"`
	Values           []string  `yaml:"with_values" bson:"values"`
}

type Condition struct {
	mgm.DefaultModel `bson:",inline"`
	Assert           string      `yaml:"assert" bson:"assert"`
	Equals           interface{} `yaml:"equals" bson:"equals"`
}

type Input struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string      `yaml:"name" bson:"name"`
	DisplayName      string      `yaml:"displayName,omitempty" bson:"displayName,omitempty"`
	Description      string      `yaml:"description,omitempty" bson:"description,omitempty"`
	Type             string      `yaml:"type" default:"string" bson:"type,omitempty"`
	Default          interface{} `yaml:"default,omitempty" bson:"default,omitempty"`
	Required         bool        `yaml:"required,omitempty" default:"false" bson:"required,omitempty"`
}

type Type struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string           `yaml:"name" bson:"name,omitempty"`
	Fields           map[string]Input `yaml:"fields" bson:"fields,omitempty"`
}

func Get(reference ResourceReference) Plugin {
	//read plugin
	provider := reference.ProviderType
	resourceType := reference.ResourceType
	mainPath := fmt.Sprintf("%s/%s/%s/%s", config.Current.Plugins.Location, provider, resourceType, MainPluginFile)
	pluginBytes, err := os.ReadFile(mainPath)
	if err != nil {
		logger.Info.Printf("File %s does not exist, skipping metadata parsing", mainPath)
		panic(errors.NotFound.WithMessage(fmt.Sprintf("Resource %s for provider %s does not exist", resourceType, provider)))
	}
	plugin := Plugin{}
	if err = yaml.Unmarshal(pluginBytes, &plugin); err != nil {
		panic(errors.ConversionError.WithMessage(err.Error()))
	}
	typePath := fmt.Sprintf("%s/%s/%s/%s", config.Current.Plugins.Location, provider, resourceType, TypePluginFile)

	if _, errExists := os.Stat(typePath); errExists == nil {
		typesBytes, errRead := os.ReadFile(typePath)
		if errRead != nil {
			panic(errors.IOError.WithMessage(errRead.Error()))
		}
		if err = yaml.Unmarshal(typesBytes, &plugin.Types); err != nil {
			panic(errors.ConversionError.WithMessage(err.Error()))
		}
	} else if errExists != nil && !os.IsNotExist(errExists) {
		panic(errors.IOError.WithMessage(errExists.Error()))
	}

	return plugin
}

func defaultMetadataPlugin(entryMap *map[string]interface{}) {
	entry := *entryMap
	if entry["monitored"] == nil || reflect.TypeOf(entry["monitored"]).Kind() != reflect.Bool {
		entry["monitored"] = true
	}
	if entry["tags"] == nil || reflect.TypeOf(entry["tags"]).Kind() != reflect.Map {
		entry["tags"] = map[string]string{}
	}
	entryMap = &entry
}

func (p *Plugin) ValidateAndCompletePluginEntry(entry map[string]interface{}) (map[string]interface{}, errors.Error) {
	defaultMetadataPlugin(&entry)
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
	notAPrimaryTypeError := func(inputType string) errors.Error {
		return errors.ValidationError.WithMessage(fmt.Sprintf("%s is not a primary type", inputType))
	}
	validatePrimitiveType := func(input Input, entry map[string]interface{}) errors.Error {
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
		return notAPrimaryTypeError(input.Type)
	}
	if err := validatePrimitiveType(i, entry); err.IsOk() {
		return errors.OK
	}
	for _, t := range types {
		if t.Name == i.Type {
			for name, input := range t.Fields {
				input.Name = name
				if err := validatePrimitiveType(input, entry[t.Name].(map[string]interface{})); reflect.DeepEqual(err, notAPrimaryTypeError(input.Type)) {
					return input.Validate(entry[t.Name].(map[string]interface{}), types)
				} else if !err.IsOk() {
					return err
				}
			}
			return errors.OK
		} else if i.Type == fmt.Sprintf("list[%s]", t.Name) {
			subEntry := entry[i.Name].([]map[string]interface{})
			for _, sub := range subEntry {
				for _, input := range t.Fields {
					if err := validatePrimitiveType(input, sub); !err.IsOk() {
						return err
					}
				}
			}
			return errors.OK
		}
	}
	return errors.OK
}

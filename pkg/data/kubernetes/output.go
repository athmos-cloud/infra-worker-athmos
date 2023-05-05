package kubernetes

import "github.com/kamva/mgm/v3"

type Output struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string      `yaml:"name" bson:"name"`
	JsonPath         string      `yaml:"jsonPath" bson:"jsonPath"`
	Value            interface{} `yaml:"value,omitempty" bson:"value,omitempty"`
}

func (output *Output) Equals(other Output) bool {
	return output.Name == other.Name &&
		output.JsonPath == other.JsonPath &&
		output.Value == other.Value
}

type OutputList []Output

func (outputList *OutputList) Equals(other OutputList) bool {
	if len(*outputList) != len(other) {
		return false
	}
	for i, output := range *outputList {
		if !output.Equals(other[i]) {
			return false
		}
	}
	return true
}

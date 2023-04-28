package kubernetes

type Output struct {
	Name     string      `bson:"name"`
	JsonPath string      `bson:"jsonPath"`
	Value    interface{} `bson:"value,omitempty"`
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

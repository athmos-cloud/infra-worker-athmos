package kubernetes

type Output struct {
	Name     string      `bson:"name"`
	JsonPath string      `bson:"jsonPath"`
	Value    interface{} `bson:"value,omitempty"`
}

type OutputList []Output

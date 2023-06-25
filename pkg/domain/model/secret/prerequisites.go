package secret

type Prerequisites struct {
	Variables []Variable `json:"variables" yaml:"variables"`
	Steps     []Command  `json:"steps" yaml:"steps"`
}

type Variable struct {
	Name    string `json:"name" yaml:"name"`
	Message string `json:"message" yaml:"message"`
}

type Command struct {
	Message string `json:"message" yaml:"message"`
	Command string `json:"command" yaml:"command"`
}

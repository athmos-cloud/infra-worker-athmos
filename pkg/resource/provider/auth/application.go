package auth

type HelmApplication struct {
	Type   Type                  `yaml:"type"`
	Secret SecretHelmApplication `yaml:"secret"`
}

type SecretHelmApplication struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Key       string `yaml:"key"`
	Value     string `yaml:"value"`
}

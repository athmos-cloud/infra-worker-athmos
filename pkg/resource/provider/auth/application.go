package auth

type SecretHelmApplication struct {
	Enabled   bool   `yaml:"enabled"`
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Key       string `yaml:"key"`
	Value     string `yaml:"value"`
}

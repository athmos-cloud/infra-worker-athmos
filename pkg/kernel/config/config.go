package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sync"
)

var Current *Config
var lock = &sync.Mutex{}

const DefaultConfigFileLocation = "config.mapstructure"

type Config struct {
	TempDir        string     `mapstructure:"tmpDir" env:"TMP_DIR"`
	RedirectionURL string     `mapstructure:"redirectionURL" env:"REDIRECTION_URL"`
	Test           Test       `mapstructure:"test" `
	Http           Http       `mapstructure:"http" `
	Queue          Queue      `mapstructure:"rabbitmq" `
	Kubernetes     Kubernetes `mapstructure:"kubernetes" `
	Mongo          Mongo      `mapstructure:"mongo" `
}

type Queue struct {
	URI      string `mapstructure:"uri"`
	Exchange string `mapstructure:"exchange"`
	Queue    string `mapstructure:"rabbitmq"`
}

type Http struct {
	Port int `mapstructure:"port" env:"PORT"`
}
type Test struct {
	Credentials CredentialsTest `mapstructure:"credentials" `
}

type CredentialsTest struct {
	GCP string `mapstructure:"resource"`
}

type Kubernetes struct {
	ConfigPath string `mapstructure:"configPath" env:"KUBECONFIG_PATH"`
}

type Mongo struct {
	Address           string `mapstructure:"address"`
	Port              int    `mapstructure:"port"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	Database          string `mapstructure:"database"`
	ProjectCollection string `mapstructure:"projectCollection"`
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if Current == nil {
		Current = &Config{}
		configPath := os.Getenv("CONFIG_FILE_LOCATION")
		if configPath == "" {
			configPath = DefaultConfigFileLocation
		}
		viper.SetConfigName(filepath.Base(configPath))
		viper.AddConfigPath(filepath.Dir(configPath))
		viper.SetConfigType("yaml")
		bindEnvs()
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(Current); err != nil {
			panic(err)
		}
	}
}

func bindEnvs() {
	if err := viper.BindEnv("redirectionURL", "REDIRECTION_URL"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("tmpDir", "TMP_DIR"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("plugins.location", "PLUGINS_LOCATION"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("kubernetes.configPath", "KUBECONFIG_PATH"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("http.port", "PORT"); err != nil {
		panic(err)
	}
}

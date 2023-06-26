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
	StaticsFileDir string     `mapstructure:"staticsFileDir" env:"STATIC_FILES_DIR"`
	RedirectionURL string     `mapstructure:"redirectionURL" env:"REDIRECTION_URL"`
	Http           Http       `mapstructure:"http"`
	Queue          Queue      `mapstructure:"rabbitmq"`
	Kubernetes     Kubernetes `mapstructure:"kubernetes"`
	Mongo          Mongo      `mapstructure:"mongo"`
}

type Queue struct {
	URI            string `mapstructure:"uri"`
	Address        string `mapstructure:"address"`
	Port           int    `mapstructure:"port"`
	Password       string `mapstructure:"password"`
	Username       string `mapstructure:"username"`
	OutcomingQueue string `mapstructure:"outcoming"`
	IncomingQueue  string `mapstructure:"incoming"`
}

type Http struct {
	Port int `mapstructure:"port" env:"PORT"`
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
	if err := viper.BindEnv("staticsFileDir", "STATIC_FILES_DIR"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("tmpDir", "TMP_DIR"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("kubernetes.configPath", "KUBECONFIG_PATH"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("http.port", "PORT"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("mongo.address", "MONGO_ADDRESS"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("rabbitmq.uri", "RABBITMQ_URI"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("rabbitmq.address", "RABBITMQ_ADDRESS"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("rabbitmq.username", "RABBITMQ_USERNAME"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("rabbitmq.port", "RABBITMQ_PORT"); err != nil {
		panic(err)
	}
}

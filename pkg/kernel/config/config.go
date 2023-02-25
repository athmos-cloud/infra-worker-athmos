package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var Current = &Config{}
var lock = &sync.Mutex{}

const DefaultConfigFileLocation = "config.yaml"

type Config struct {
	TempDir   string    `yaml:"tmp_dir" env:"TMP_DIR" env-default:"/tmp/infra-worker"`
	Runtime   string    `yaml:"runtime" env:"RUNTIME" env-default:"dagger"`
	Terraform Terraform `yaml:"terraform" prefix:"TERRAFORM_"`
	Minio     Minio     `yaml:"minio" prefix:"MINIO_"`
	Mongo     Mongo     `yaml:"mongo" prefix:"MONGO_"`
	Postgres  Postgres  `yaml:"postgres" prefix:"POSTGRES_"`
}

type Terraform struct {
	Image        Image  `yaml:"image" prefix:"IMAGE_"`
	BucketPlugin string `yaml:"bucket_plugin" env:"BUCKET_PLUGIN" env-default:"terraform"`
}

type Image struct {
	Name string `yaml:"name" env:"IMAGE_NAME" env-default:"hashicorp/terraform:1.3.9"`
	Tag  string `yaml:"tag" env:"IMAGE_TAG" env-default:"1.3.9"`
}

type Minio struct {
	Address         string `yaml:"address" env:"ADDRESS" env-default:"minio"`
	AccessKeyID     string `yaml:"access_key_id" env:"ACCESS_KEY_ID"`
	SecretAccessKey string `yaml:"secret_access_key" env:"SECRET_ACCESS_KEY"`
	Token           string `yaml:"token" env:"TOKEN" env-default:""`
	Region          string `yaml:"region" env:"REGION" env-default:""`
	UseSSL          bool   `yaml:"use_ssl" env:"USE_SSL" env-default:"false"`
}

type Mongo struct {
	Address  string `yaml:"address" env:"ADDRESS" env-default:"mongo"`
	Port     int    `yaml:"port" env:"PORT" env-default:"27017"`
	Username string `yaml:"username" env:"USERNAME" env-default:"root"`
	Password string `yaml:"password" env:"PASSWORD"`
	Database string `yaml:"database" env:"DATABASE" env-default:"plugin-db"`
}

type Postgres struct {
	Address  string `yaml:"host" env:"ADDRESS" env-default:"postgres"`
	Port     int    `yaml:"port" env:"PORT" env-default:"5432"`
	Username string `yaml:"username" env:"USERNAME" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD"`
	Database string `yaml:"database" env:"DATABASE" env-default:"plugin-db"`
	SSLMode  string `yaml:"ssl_mode" env:"SSL_MODE" env-default:"disable"`
}

func Get() *Config {
	if Current == nil {
		lock.Lock()
		defer lock.Unlock()
		readEnv(Current)
		readFile(Current)
	}
	return Current
}

func readFile(cfg *Config) {
	configFile := os.Getenv("CONFIG_FILE_LOCATION")
	if configFile == "" {
		configFile = DefaultConfigFileLocation
	}
	f, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err1 := f.Close()
		if err1 != nil {
			panic(err1)
		}
	}(f)

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		panic(err)
	}
}

func readEnv(cfg *Config) {
	err := cleanenv.ReadEnv(&Current)
	if err != nil {
		return
	}
}

package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tkanos/gonfig"
	"os"
	"sync"
)

var Current *Config
var lock = &sync.Mutex{}

const DefaultConfigFileLocation = "config.yaml"

type Config struct {
	TempDir    string     `yaml:"tmp_dir" env:"TMP_DIR" env-default:"/tmp/infra-worker"`
	Kubernetes Kubernetes `yaml:"kubernetes" prefix:"KUBERNETES_"`
	Plugins    Plugins    `yaml:"plugins" prefix:"PLUGINS_"`
	Mongo      Mongo      `yaml:"mongo" prefix:"MONGO_"`
	Postgres   Postgres   `yaml:"postgres" prefix:"POSTGRES_"`
}

type Kubernetes struct {
	ConfigPath string `yaml:"configPath" env:"KUBECONFIG_PATH" env-default:"~/.kube/config"`
	Helm       Helm   `yaml:"helm" prefix:"HELM_"`
}

type Helm struct {
	Debug bool `yaml:"debug" env:"DEBUG" env-default:"false"`
}

type Plugins struct {
	Location   string            `yaml:"location" env:"LOCATION" env-default:"/plugins"`
	Crossplane CrossplanePlugins `yaml:"crossplane" prefix:"CROSSPLANE_"`
}

type CrossplanePlugins struct {
	Registry ArtifactRegistry `yaml:"artifact-registry" prefix:"ARTIFACT_REGISTRY_"`
	GCP      ProviderPlugins  `yaml:"gcp"`
}

type ArtifactRegistry struct {
	Address  string `yaml:"address" env:"ADDRESS"`
	Username string `yaml:"username" env:"USERNAME"`
	Password string `yaml:"password" env:"PASSWORD"`
}

type ProviderPlugins struct {
	Firewall ProviderPluginItem `yaml:"firewall"`
	Network  ProviderPluginItem `yaml:"network"`
	Provider ProviderPluginItem `yaml:"provider"`
	Subnet   ProviderPluginItem `yaml:"subnetwork"`
	VM       ProviderPluginItem `yaml:"vm"`
	VPC      ProviderPluginItem `yaml:"vpc"`
}

type ProviderPluginItem struct {
	Chart   string `yaml:"chart"`
	Version string `yaml:"version"`
}

type Mongo struct {
	Address           string `yaml:"address" env:"ADDRESS" env-default:"mongo"`
	Port              int    `yaml:"port" env:"PORT" env-default:"27017"`
	Username          string `yaml:"username" env:"USERNAME" env-default:"root"`
	Password          string `yaml:"password" env:"PASSWORD"`
	Database          string `yaml:"database" env:"DATABASE" env-default:"plugin-db"`
	ProjectCollection string `yaml:"project_collection" env:"PROJECT_COLLECTION" env-default:"projects"`
}

type Postgres struct {
	Address  string `yaml:"host" env:"ADDRESS" env-default:"postgres"`
	Port     int    `yaml:"port" env:"PORT" env-default:"5432"`
	Username string `yaml:"username" env:"USERNAME" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD"`
	Database string `yaml:"database" env:"DATABASE" env-default:"plugin-db"`
	SSLMode  string `yaml:"ssl_mode" env:"SSL_MODE" env-default:"disable"`
}

func init() {
	Current = &Config{}
	readEnv()
	readFile()
}

func readFile() {
	configFile := os.Getenv("CONFIG_FILE_LOCATION")
	if configFile == "" {
		configFile = DefaultConfigFileLocation
	}
	err := gonfig.GetConf(configFile, Current)
	if err != nil {
		panic(err)
	}
}

func readEnv() {
	err := cleanenv.ReadEnv(Current)
	if err != nil {
		return
	}
}

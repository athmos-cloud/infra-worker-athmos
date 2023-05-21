package share

type EnvType string

const (
	EnvTypeDevelopment EnvType = "development"
	EnvTypeProduction  EnvType = "production"
	EnvTypeTest        EnvType = "test"
)

const (
	EnvTypeEnvironmentVariable = "ENV_TYPE"
)

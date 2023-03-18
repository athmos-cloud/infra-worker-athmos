package auth

type AuthType string

const (
	SecretAuthType  AuthType = "Secret"
	SecretVaultType AuthType = "SecretVault"
)

type Auth struct {
	AuthType    AuthType    `bson:"authType"`
	SecretAuth  SecretAuth  `bson:"secretAuth"`
	SecretVault SecretVault `bson:"secretVault"`
}

type SecretAuth struct {
	SecretName string `bson:"secretName"`
	SecretKey  string `bson:"secretKey"`
	Namespace  string `bson:"namespace"`
}

type SecretVault struct{}

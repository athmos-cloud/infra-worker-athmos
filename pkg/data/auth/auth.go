package domain

type Type string

const (
	AuthTypeSecret Type = "secret"
	AuthTypeVault  Type = "vault"
)

type Auth struct {
	AuthType    Type        `bson:"authType"`
	SecretAuth  SecretAuth  `bson:"secretAuth"`
	SecretVault SecretVault `bson:"secretVault"`
}

type SecretAuth struct {
	SecretName string `bson:"secretName"`
	SecretKey  string `bson:"secretKey"`
	Namespace  string `bson:"namespace"`
}

type SecretVault struct{}

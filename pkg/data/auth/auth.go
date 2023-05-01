package auth

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

const (
	DefaultSecretKey = "key.json"
)

type Type string

const (
	AuthTypeSecret Type = "secret"
	AuthTypeVault  Type = "vault"
)

func AuthType(str string) (Type, errors.Error) {
	switch str {
	case "secret":
		return AuthTypeSecret, errors.OK
	case "vault":
		return AuthTypeVault, errors.OK
	default:
		return "", errors.InvalidArgument.WithMessage(fmt.Sprintf("Auth type %s not supported", str))
	}
}

type Auth struct {
	Name        string      `bson:"name"`
	Description string      `bson:"description"`
	AuthType    Type        `bson:"authType" plugin:"authType"`
	SecretAuth  SecretAuth  `bson:"secretAuth" plugin:"secret"`
	SecretVault SecretVault `bson:"secretVault" plugin:"vault"`
}

type AuthList map[string]Auth

func (a *Auth) Equals(auth Auth) bool {
	return a.AuthType == auth.AuthType &&
		a.SecretAuth.Equals(auth.SecretAuth) &&
		a.SecretVault.Equals(auth.SecretVault)
}

type SecretAuth struct {
	SecretName string `bson:"secretName" plugin:"name"`
	SecretKey  string `bson:"secretKey" plugin:"key"`
	Namespace  string `bson:"namespace" plugin:"namespace"`
}

func (a *SecretAuth) Equals(auth SecretAuth) bool {
	return a.SecretName == auth.SecretName &&
		a.SecretKey == auth.SecretKey &&
		a.Namespace == auth.Namespace
}

type SecretVault struct{}

func (v *SecretVault) Equals(vault SecretVault) bool {
	return true
}

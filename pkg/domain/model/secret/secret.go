package secret

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/kamva/mgm/v3"
)

const (
	NameLabelKey        = "name.secret"
	DescriptionLabelKey = "description.secret"
)

type Secret struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string         `bson:"name"`
	Description      string         `bson:"description,omitempty"`
	ProviderType     types.Provider `bson:"provider_type"`
	Kubernetes       Kubernetes     `bson:"secret_auth,omitempty"`
	Prerequisites    Prerequisites  `bson:"prerequisites,omitempty"`
}

type List map[string]Secret

func NewSecret(name string, description string, secretAuth Kubernetes, forProvider types.Provider) *Secret {
	return &Secret{
		Name:         name,
		ProviderType: forProvider,
		Description:  description,
		Kubernetes:   secretAuth,
	}
}

type Kubernetes struct {
	mgm.DefaultModel `bson:",inline"`
	SecretName       string `bson:"secretName" plugin:"name"`
	SecretKey        string `bson:"secretKey" plugin:"key"`
	Namespace        string `bson:"namespace" plugin:"namespace"`
}

func NewKubernetesSecret(secretName string, secretKey string, namespace string) Kubernetes {
	return Kubernetes{
		SecretName: secretName,
		SecretKey:  secretKey,
		Namespace:  namespace,
	}
}

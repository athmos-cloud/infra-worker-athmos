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
	ProviderType     types.Provider `bson:"providerType"`
	Kubernetes       Kubernetes     `bson:"secretAuth,omitempty"`
}

type List map[string]Secret

func (s *Secret) Equals(other Secret) bool {
	return s.Name == other.Name && s.Description == other.Description && s.Kubernetes.Equals(other.Kubernetes)
}

func NewSecret(name string, description string, secretAuth Kubernetes) *Secret {
	return &Secret{
		Name:        name,
		Description: description,
		Kubernetes:  secretAuth,
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

func (s *Kubernetes) Equals(auth Kubernetes) bool {
	return s.SecretName == auth.SecretName &&
		s.SecretKey == auth.SecretKey &&
		s.Namespace == auth.Namespace
}

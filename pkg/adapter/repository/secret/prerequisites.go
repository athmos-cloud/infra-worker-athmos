package secret

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	secretRepos "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/secret"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type yamlPrerequisitesRepository struct{}

func NewYamlPrerequisitesRepository() secretRepos.PrerequisitesRepository {
	return &yamlPrerequisitesRepository{}
}

func (pr *yamlPrerequisitesRepository) Find(prerequisiteSecret *secret.Secret) errors.Error {
	file := strings.TrimSuffix(config.Current.StaticsFileDir, "/")
	switch prerequisiteSecret.ProviderType {
	case types.ProviderGCP:
		file += "/gcp/prerequisites.yaml"
	case types.ProviderAWS:
		file += "/aws/prerequisites.yaml"
	case types.ProviderAZURE:
		file += "/azure/prerequisites.yaml"
	}
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return errors.InternalError.WithMessage(fmt.Sprintf("Unable to read prerequisites file : %s", file))
	}
	prerequisites := secret.Prerequisites{}
	if errUnmarshall := yaml.Unmarshal(yamlFile, &prerequisites); errUnmarshall != nil {
		return errors.InternalError.WithMessage(errUnmarshall.Error())
	}
	prerequisiteSecret.Prerequisites = prerequisites
	return errors.OK
}

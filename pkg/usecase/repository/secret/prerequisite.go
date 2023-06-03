package secret

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type PrerequisitesRepository interface {
	Find(*secret.Secret) errors.Error
}

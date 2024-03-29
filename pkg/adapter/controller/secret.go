package controller

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	error2 "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

type Secret interface {
	ListProjectSecret(ctx context.Context)
	GetSecret(ctx context.Context)
	CreateSecret(ctx context.Context)
	UpdateSecret(ctx context.Context)
	DeleteSecret(ctx context.Context)
}

type secretController struct {
	secretUseCase usecase.Secret
	secretPort    output.SecretPort
}

func NewSecretController(secretUseCase usecase.Secret, secretPort output.SecretPort) Secret {
	return &secretController{
		secretUseCase: secretUseCase,
		secretPort:    secretPort,
	}
}

func (s *secretController) ListProjectSecret(ctx context.Context) {
	if err := validator.ListProjectSecret(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	secrets := &[]secret.Secret{}
	if err := s.secretUseCase.List(ctx, secrets); !err.IsOk() {
		error2.RaiseError(ctx, err)
	} else {
		s.secretPort.RenderAll(ctx, secrets)
	}
}

func (s *secretController) GetSecret(ctx context.Context) {
	if err := validator.GetSecret(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	secretAuth := &secret.Secret{}
	if err := s.secretUseCase.Get(ctx, secretAuth); !err.IsOk() {
		error2.RaiseError(ctx, err)
	} else {
		s.secretPort.Render(ctx, secretAuth)
	}
}

func (s *secretController) CreateSecret(ctx context.Context) {
	if err := validator.CreateSecret(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	secretAuth := &secret.Secret{}
	if err := s.secretUseCase.Create(ctx, secretAuth); !err.IsOk() {
		error2.RaiseError(ctx, err)
	} else {
		s.secretPort.RenderCreate(ctx, secretAuth)
	}
}

func (s *secretController) UpdateSecret(ctx context.Context) {
	if err := validator.UpdateSecret(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	secretAuth := &secret.Secret{}
	if err := s.secretUseCase.Update(ctx, secretAuth); !err.IsOk() {
		error2.RaiseError(ctx, err)
	} else {
		s.secretPort.Render(ctx, secretAuth)
	}
}

func (s *secretController) DeleteSecret(ctx context.Context) {
	if err := validator.DeleteSecret(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	if err := s.secretUseCase.Delete(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	} else {
		s.secretPort.RenderDelete(ctx)
	}
}

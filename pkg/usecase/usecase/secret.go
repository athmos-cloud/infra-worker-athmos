package usecase

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	secretRepos "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	secret2 "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/secret"
)

const (
	defaultSecretKey = "credentials.json"
)

type Secret interface {
	Get(context.Context, *secret.Secret) errors.Error
	List(context.Context, *[]secret.Secret) errors.Error
	Create(context.Context, *secret.Secret) errors.Error
	Update(context.Context, *secret.Secret) errors.Error
	Delete(context.Context) errors.Error
}

type secretUseCase struct {
	secretRepo           secret2.Secret
	prerequisitesRepo    secret2.PrerequisitesRepository
	kubernetesSecretRepo secret2.KubernetesSecret
}

func NewSecretUseCase(secretRepo secret2.Secret, kubernetesRepo secret2.KubernetesSecret) Secret {
	return &secretUseCase{secretRepo: secretRepo, prerequisitesRepo: secretRepos.NewYamlPrerequisitesRepository(), kubernetesSecretRepo: kubernetesRepo}
}

func (suc *secretUseCase) Get(ctx context.Context, secretAuth *secret.Secret) errors.Error {
	req := ctx.Value(context.RequestKey).(dto.GetSecretRequest)
	foundSecret, err := suc.secretRepo.Find(ctx, option.Option{
		Value: secret2.GetSecretByProjectIdAndName{
			ProjectId: req.ProjectID,
			Name:      req.Name,
		},
	})
	if !err.IsOk() {
		return err
	}
	*secretAuth = *foundSecret
	return errors.OK
}

func (suc *secretUseCase) List(ctx context.Context, secretAuths *[]secret.Secret) errors.Error {
	foundSecrets, err := suc.secretRepo.FindAll(ctx, option.Option{
		Value: secret2.GetSecretInProject{
			ProjectId: ctx.Value(context.ProjectIDKey).(string),
		},
	})
	if !err.IsOk() {
		return err
	}
	*secretAuths = *foundSecrets
	return errors.OK
}

func (suc *secretUseCase) Create(ctx context.Context, secretAuth *secret.Secret) errors.Error {
	req := ctx.Value(context.RequestKey).(dto.CreateSecretRequest)
	secretName := IdFromName(req.Name)

	createdSecret, err := suc.kubernetesSecretRepo.Create(ctx, option.Option{
		Value: secret2.CreateKubernetesSecretRequest{
			ProjectID:   req.ProjectID,
			SecretName:  secretName,
			SecretKey:   defaultSecretKey,
			SecretValue: []byte(req.Value),
		},
	})
	if !err.IsOk() {
		return err
	}
	*secretAuth = *secret.NewSecret(req.Name, req.Description, *createdSecret, req.ForProvider)
	if errRepo := suc.secretRepo.Create(ctx, secretAuth); !errRepo.IsOk() {
		return errRepo
	}
	if errPrerequisites := suc.prerequisitesRepo.Find(secretAuth); !errPrerequisites.IsOk() {
		return errPrerequisites
	}
	return errors.NoContent
}

func (suc *secretUseCase) Update(ctx context.Context, secretAuth *secret.Secret) errors.Error {
	req := ctx.Value(context.RequestKey).(dto.UpdateSecretRequest)
	projectID := ctx.Value(context.ProjectIDKey).(string)
	curSecret, err := suc.secretRepo.Find(ctx, option.Option{
		Value: secret2.GetSecretByProjectIdAndName{
			ProjectId: projectID,
			Name:      req.Name,
		},
	})
	if !err.IsOk() {
		return err
	}
	if req.Value != nil {
		if errKube := suc.kubernetesSecretRepo.Update(ctx, option.Option{
			Value: secret2.UpdateKubernetesSecretRequest{
				ProjectID:   ctx.Value(context.ProjectIDKey).(string),
				SecretName:  curSecret.Kubernetes.SecretName,
				SecretKey:   defaultSecretKey,
				SecretValue: req.Value[:],
			},
		}); !errKube.IsOk() {
			return errKube
		}
	}
	*secretAuth = *curSecret
	if req.Description == "" {
		return errors.NoContent
	}
	secretAuth.Description = req.Description
	if errRepo := suc.secretRepo.Update(ctx, secretAuth); !errRepo.IsOk() {
		return errRepo
	}
	return errors.NoContent
}

func (suc *secretUseCase) Delete(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(dto.DeleteSecretRequest)
	projectID := ctx.Value(context.ProjectIDKey).(string)
	curSecret, err := suc.secretRepo.Find(ctx, option.Option{
		Value: secret2.GetSecretByProjectIdAndName{
			ProjectId: projectID,
			Name:      req.Name,
		},
	})
	if !err.IsOk() {
		return err
	}
	if errKube := suc.kubernetesSecretRepo.Delete(ctx, option.Option{
		Value: secret2.DeleteKubernetesSecretRequest{
			ProjectID:  projectID,
			SecretName: curSecret.Kubernetes.SecretName,
		},
	}); !errKube.IsOk() {
		return errKube
	}

	if errDelete := suc.secretRepo.Delete(ctx, curSecret); !errDelete.IsOk() {
		return errDelete
	}
	return errors.NoContent
}

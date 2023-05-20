package usecase

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	arepo "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
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
	secretRepo           repository.Secret
	kubernetesSecretRepo repository.KubernetesSecret
}

func NewSecretUseCase(secretRepo repository.Secret, kubernetesRepo repository.KubernetesSecret) Secret {
	return &secretUseCase{secretRepo: secretRepo, kubernetesSecretRepo: kubernetesRepo}
}

func (suc *secretUseCase) Get(ctx context.Context, secretAuth *secret.Secret) errors.Error {
	req := ctx.Value(context.RequestKey).(dto.GetSecretRequest)
	foundSecret, err := suc.secretRepo.Find(ctx, option.Option{
		Value: arepo.GetSecretByProjectIdAndName{
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
		Value: arepo.GetSecretInProject{
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
	secretName := fmt.Sprintf("%s-%s", req.Name, utils.RandomString(5))
	createdSecret, err := suc.kubernetesSecretRepo.Create(ctx, option.Option{
		Value: arepo.CreateKubernetesSecretRequest{
			ProjectID:   req.ProjectID,
			SecretName:  secretName,
			SecretKey:   defaultSecretKey,
			SecretValue: req.Value[:],
		},
	})
	if !err.IsOk() {
		return err
	}
	*secretAuth = *secret.NewSecret(req.Name, req.Description, *createdSecret)
	if errRepo := suc.secretRepo.Create(ctx, secretAuth); !errRepo.IsOk() {
		return errRepo
	}
	return errors.NoContent
}

func (suc *secretUseCase) Update(ctx context.Context, secretAuth *secret.Secret) errors.Error {
	req := ctx.Value(context.RequestKey).(dto.UpdateSecretRequest)
	projectID := ctx.Value(context.ProjectIDKey).(string)
	curSecret, err := suc.secretRepo.Find(ctx, option.Option{
		Value: arepo.GetSecretByProjectIdAndName{
			ProjectId: projectID,
			Name:      req.Name,
		},
	})
	if !err.IsOk() {
		return err
	}
	if req.Value != nil {
		if errKube := suc.kubernetesSecretRepo.Update(ctx, option.Option{
			Value: arepo.UpdateKubernetesSecretRequest{
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
	if err := suc.kubernetesSecretRepo.Delete(ctx, option.Option{
		Value: arepo.DeleteKubernetesSecretRequest{
			ProjectID:  projectID,
			SecretName: req.Name,
		},
	}); !err.IsOk() {
		return err
	}
	curSecret, err := suc.secretRepo.Find(ctx, option.Option{
		Value: arepo.GetSecretByProjectIdAndName{
			ProjectId: projectID,
			Name:      req.Name,
		},
	})
	if !err.IsOk() {
		return err
	}
	if errDelete := suc.secretRepo.Delete(ctx, curSecret); !errDelete.IsOk() {
		return errDelete
	}
	return errors.NoContent
}

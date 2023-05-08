package secret

import (
	"context"
	"fmt"
	project2 "github.com/athmos-cloud/infra-worker-athmos/pkg/application/project"
	kubernetesDAO "github.com/athmos-cloud/infra-worker-athmos/pkg/dao/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/auth"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	projectRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/project"
)

const randomSecretStringLength = 5

type Service struct {
	KubernetesDAO     *kubernetesDAO.DAO
	ProjectRepository *projectRepo.Repository
}

func (service *Service) CreateSecret(ctx context.Context, request CreateSecretRequest) {
	// Get the project
	project := service.ProjectRepository.Get(ctx, option.Option{
		Value: project2.GetProjectByIDRequest{
			ProjectID: request.ProjectID,
		},
	})
	currentProject := project.(project2.GetProjectByIDResponse).Payload
	logger.Info.Printf("Project FirewallID : %s", request.ProjectID)

	if _, ok := currentProject.Authentications[request.Name]; ok {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("Secret %s already exists", request.Name)))
	}
	secretName := fmt.Sprintf("%s-%s", request.Name, utils.RandomString(randomSecretStringLength))
	// Create the secret into Kubernetes
	service.KubernetesDAO.Create(ctx, option.Option{
		Value: kubernetesDAO.CreateSecretRequest{
			Name:      secretName,
			Namespace: currentProject.Namespace,
			Key:       auth.DefaultSecretKey,
			Data:      []byte(request.Data),
		},
	})
	// Insert secret into the project
	currentProject.Authentications[request.Name] = auth.Auth{
		Name:        request.Name,
		Description: request.Description,
		AuthType:    auth.AuthTypeSecret,
		SecretAuth: auth.SecretAuth{
			SecretName: secretName,
			SecretKey:  auth.DefaultSecretKey,
			Namespace:  currentProject.Namespace,
		},
	}

	// Persist the project into the database
	service.ProjectRepository.Update(ctx, option.Option{
		Value: projectRepo.UpdateProjectRequest{
			ProjectID:      request.ProjectID,
			UpdatedProject: currentProject,
		},
	})

}

func (service *Service) GetSecret(ctx context.Context, request GetSecretRequest) GetSecretResponse {
	project := service.ProjectRepository.Get(ctx, option.Option{
		Value: project2.GetProjectByIDRequest{
			ProjectID: request.ProjectID,
		},
	})
	currentProject := project.(project2.GetProjectByIDResponse).Payload
	currentSecret, exists := currentProject.Authentications[request.Name]

	if !exists {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("Secret %s not found", request.Name)))
	}
	return GetSecretResponse{
		Name:        currentSecret.Name,
		Description: currentSecret.Description,
		References:  currentSecret.SecretAuth,
	}
}

func (service *Service) ListSecret(ctx context.Context, request ListSecretRequest) ListSecretResponse {
	project := service.ProjectRepository.Get(ctx, option.Option{
		Value: project2.GetProjectByIDRequest{
			ProjectID: request.ProjectID,
		},
	})
	currentProject := project.(project2.GetProjectByIDResponse).Payload

	var secrets ListSecretResponse
	for _, s := range currentProject.Authentications {
		secrets = append(secrets, GetSecretResponse{
			Name:        s.Name,
			Description: s.Description,
		})
	}
	return secrets
}

func (service *Service) UpdateSecret(ctx context.Context, request UpdateSecretRequest) {
	project := service.ProjectRepository.Get(ctx, option.Option{
		Value: project2.GetProjectByIDRequest{
			ProjectID: request.ProjectID,
		},
	})

	currentProject := project.(project2.GetProjectByIDResponse).Payload
	if _, exists := currentProject.Authentications[request.Name]; !exists {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("Cannot update secret %s", request.Name)))
	}

	secretName := fmt.Sprintf("%s-%s", request.Name, utils.RandomString(randomSecretStringLength))
	currentProject.Authentications[request.Name] = auth.Auth{
		Name:        request.Name,
		Description: request.Description,
		AuthType:    auth.AuthTypeSecret,
		SecretAuth: auth.SecretAuth{
			SecretName: secretName,
			SecretKey:  auth.DefaultSecretKey,
			Namespace:  currentProject.Namespace,
		},
	}

	service.KubernetesDAO.Update(ctx, option.Option{
		Value: kubernetesDAO.UpdateSecretRequest{
			Name:      secretName,
			Namespace: currentProject.Namespace,
			Key:       auth.DefaultSecretKey,
			Data:      []byte(request.Data),
		},
	})

	service.ProjectRepository.Update(ctx, option.Option{
		Value: projectRepo.UpdateProjectRequest{
			ProjectID:      request.ProjectID,
			UpdatedProject: currentProject,
		},
	})
}

func (service *Service) DeleteSecret(ctx context.Context, request DeleteSecretRequest) {
	project := service.ProjectRepository.Get(ctx, option.Option{
		Value: project2.GetProjectByIDRequest{
			ProjectID: request.ProjectID,
		},
	})
	currentProject := project.(project2.GetProjectByIDResponse).Payload

	secretToDelete, ok := currentProject.Authentications[request.Name]
	if !ok {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("Invalid secret %s", request.Name)))
	}

	delete(currentProject.Authentications, request.Name)

	service.KubernetesDAO.Delete(ctx, option.Option{
		Value: kubernetesDAO.DeleteSecretRequest{
			Namespace: secretToDelete.SecretAuth.Namespace,
			Name:      request.Name,
		},
	})

	service.ProjectRepository.Update(ctx, option.Option{
		Value: projectRepo.UpdateProjectRequest{
			ProjectID:      request.ProjectID,
			UpdatedProject: currentProject,
		},
	})
}

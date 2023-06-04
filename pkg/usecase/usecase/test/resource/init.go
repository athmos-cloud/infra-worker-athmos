package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	repositoryAdapter "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	secretModel "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	domainTypes "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	secret2 "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	ProjectTestName       = "test"
	ProjectTestOwner      = "test"
	SecretTestName        = "test"
	SecretTestDescription = "Some test secret"
	ProviderIDKey         = "provider_id"
	NetworkIDKey          = "network_id"
	SubnetworkIDKey       = "subnetwork_id"
)

type TestResource struct {
	ProjectRepo          repository.Project
	SecretRepo           secret2.Secret
	KubernetesSecretRepo secret2.KubernetesSecret
}

func NewTest(ctx context.Context, t *testing.T) *TestResource {
	projectRepo := repositoryAdapter.NewProjectRepository()
	secretRepo := secret.NewSecretRepository()
	kubeSecretRepo := secret.NewKubernetesRepository()

	project := model.NewProject(ProjectTestName, ProjectTestOwner)
	err := projectRepo.Create(ctx, project)
	require.True(t, err.IsOk())
	ctx.Set(context.ProjectIDKey, project.ID.Hex())

	createdSecret, err := kubeSecretRepo.Create(ctx, option.Option{
		Value: secret2.CreateKubernetesSecretRequest{
			ProjectID:   project.ID.Hex(),
			SecretName:  "test",
			SecretKey:   "key.json",
			SecretValue: []byte("{\"test\":\"test\"}"),
		},
	})
	require.True(t, err.IsOk())
	secretAuth := secretModel.NewSecret(SecretTestName, SecretTestDescription, *createdSecret, domainTypes.ProviderGCP)
	err = secretRepo.Create(ctx, secretAuth)
	require.True(t, err.IsOk())
	ctx.Set(test.TestNamespaceContextKey, project.Namespace)
	ctx.Set(test.TestSecretContextKey, *secretAuth)

	return &TestResource{
		ProjectRepo:          projectRepo,
		SecretRepo:           secretRepo,
		KubernetesSecretRepo: kubeSecretRepo,
	}
}

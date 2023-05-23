package gcp

import (
	repository2 "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_providerUseCase_Create(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	projectRepo := repository2.NewProjectRepository()
	secretRepo := secret.NewSecretRepository()
	gcpRepo := gcp.NewRepository()
	uc := usecase.NewProviderUseCase(projectRepo, secretRepo, gcpRepo, nil, nil)
	ctx := test.NewContext()
	project := model.NewProject("test", "test")
	projectRepo.Create(project)
	t.Run("Create a valid provider", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a provider with a non-existing secret should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a provider with an already existing name should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_Delete(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Delete a valid provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing provider should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a provider with children should fail", func(t *testing.T) {

	})
	t.Run("Delete cascade a provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_Get(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Get a valid provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing provider should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_GetRecursively(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("GetRecursively a valid provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_List(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("List providers should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("List providers in a non-existing project should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_Update(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Update a valid provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Update a non-existing provider should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

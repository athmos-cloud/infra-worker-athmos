package test

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	secretRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	secretModel "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

const (
	defaultSecretKey = "credentials.json"
)

func Test_secretUseCase_Create(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	pu := usecase.NewProjectUseCase(repository.NewProjectRepository())
	curProject := model.NewProject("test", "test")
	ctx := NewContext()
	err := pu.Create(ctx, curProject)
	ctx = ctx.WithValue(share.ProjectIDKey, curProject.ID.Hex())
	require.True(t, err.IsOk())

	t.Run("Should successfully create a secret", func(t *testing.T) {
		secretData := "test"
		secretName := "test-secret"
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())
		ctx = ctx.WithValue(share.RequestContextKey, dto.CreateSecretRequest{
			ProjectID:   curProject.ID.Hex(),
			Name:        secretName,
			Description: "A test secret",
			Value:       []byte(secretData),
		})
		secret := &secretModel.Secret{}
		errCreate := su.Create(ctx, secret)
		assert.True(t, errCreate.IsOk())
		ctx = ctx.WithValue(share.RequestContextKey, dto.GetSecretRequest{
			ProjectID: curProject.ID.Hex(),
			Name:      secretName,
		})

		gotSecret := &secretModel.Secret{}
		errGet := su.Get(ctx, gotSecret)
		assert.True(t, errGet.IsOk())
		kubeSecret := &corev1.Secret{}
		errKube := kubernetes.Client().Get(ctx, types.NamespacedName{
			Name:      gotSecret.Kubernetes.SecretName,
			Namespace: curProject.Namespace,
		}, kubeSecret)

		assert.Nil(t, errKube)
		assert.Equal(t, secretData, string(kubeSecret.Data[defaultSecretKey]))
	})
	t.Run("Should fail to create a secret when a secret with same name already exists in project", func(t *testing.T) {
		secretData := "test"
		secretName := "test-secret-2"
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())
		ctx = ctx.WithValue(share.RequestContextKey, dto.CreateSecretRequest{
			ProjectID:   curProject.ID.Hex(),
			Name:        secretName,
			Description: "A test secret",
			Value:       []byte(secretData),
		})
		secret := &secretModel.Secret{}
		errCreate := su.Create(ctx, secret)
		assert.True(t, errCreate.IsOk())
		errRecreate := su.Create(ctx, secret)
		assert.Equal(t, errors.Conflict.Code, errRecreate.Code)
	})
}

func Test_secretUseCase_Delete(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	pu := usecase.NewProjectUseCase(repository.NewProjectRepository())
	curProject := model.NewProject("test", "test")
	ctx := NewContext()
	err := pu.Create(ctx, curProject)
	require.True(t, err.IsOk())
	ctx = ctx.WithValue(share.ProjectIDKey, curProject.ID.Hex())

	t.Run("Should successfully delete a secret", func(t *testing.T) {
		secretName := "test-secret"
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())
		secret := &secretModel.Secret{}
		ctx = ctx.WithValue(share.RequestContextKey, dto.CreateSecretRequest{
			ProjectID:   curProject.ID.Hex(),
			Name:        secretName,
			Description: "A test secret",
			Value:       []byte("test"),
		})
		errCreate := su.Create(ctx, secret)
		assert.True(t, errCreate.IsOk())
		ctx = ctx.WithValue(share.RequestContextKey, dto.DeleteSecretRequest{
			Name: secret.Kubernetes.SecretName,
		})
		errGet := su.Delete(ctx)
		assert.Equal(t, errors.NotFound.Code, errGet.Code)
		kubeSecret := &corev1.Secret{}
		errKube := kubernetes.Client().Get(ctx, types.NamespacedName{
			Name:      secret.Kubernetes.SecretName,
			Namespace: curProject.Namespace,
		}, kubeSecret)
		assert.NotNil(t, errKube)
		assert.True(t, k8serrors.IsNotFound(errKube))
	})

	t.Run("Should fail to delete a secret when a secret with same name already exists in project", func(t *testing.T) {
		secretName := "test-secret"
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())
		ctx = ctx.WithValue(share.RequestContextKey, dto.DeleteSecretRequest{
			Name: secretName,
		})
		errGet := su.Delete(ctx)
		assert.Equal(t, errors.NotFound.Code, errGet.Code)
	})
}

func Test_secretUseCase_Get(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()

	pu := usecase.NewProjectUseCase(repository.NewProjectRepository())
	curProject := model.NewProject("test", "test")
	ctx := NewContext()
	err := pu.Create(ctx, curProject)
	require.True(t, err.IsOk())
	ctx = ctx.WithValue(share.ProjectIDKey, curProject.ID.Hex())

	t.Run("Get existing secret", func(t *testing.T) {
		secretData := "test"
		secretName := "test-secret"
		description := "A test secret"
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())
		ctx = ctx.WithValue(share.RequestContextKey, dto.CreateSecretRequest{
			ProjectID:   curProject.ID.Hex(),
			Name:        secretName,
			Description: description,
			Value:       []byte(secretData),
		})
		secret := &secretModel.Secret{}
		errCreate := su.Create(ctx, secret)
		assert.True(t, errCreate.IsOk())
		ctx = ctx.WithValue(share.RequestContextKey, dto.GetSecretRequest{
			ProjectID: curProject.ID.Hex(),
			Name:      secretName,
		})
		errGet := su.Get(ctx, secret)
		assert.True(t, errGet.IsOk())
		assert.Equal(t, secret.Name, secretName)
		assert.Equal(t, secret.Description, description)

	})

	t.Run("Get non-existing secret should return NotFound error", func(t *testing.T) {
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())
		ctx = ctx.WithValue(share.RequestContextKey, dto.GetSecretRequest{
			ProjectID: curProject.ID.Hex(),
			Name:      "non-existing-secret",
		})
		secret := &secretModel.Secret{}
		errGet := su.Get(ctx, secret)
		assert.Equal(t, errors.NotFound.Code, errGet.Code)
	})

}

func Test_secretUseCase_List(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()

	pu := usecase.NewProjectUseCase(repository.NewProjectRepository())
	curProject := model.NewProject("test", "test")
	ctx := NewContext()
	err := pu.Create(ctx, curProject)
	require.True(t, err.IsOk())
	ctx = ctx.WithValue(share.ProjectIDKey, curProject.ID.Hex())

	t.Run("List secrets in a project without secrets", func(t *testing.T) {
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())
		secrets := &[]secretModel.Secret{}
		errList := su.List(ctx, secrets)
		assert.True(t, errList.IsOk())
		assert.Equal(t, 0, len(*secrets))
	})

	t.Run("List secrets", func(t *testing.T) {
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())

		// CreateNetwork a first secret
		ctx = ctx.WithValue(share.RequestContextKey, dto.CreateSecretRequest{
			ProjectID:   curProject.ID.Hex(),
			Name:        "test-secret-1",
			Description: "A test secret 1",
			Value:       []byte("test1"),
		})
		secret1 := &secretModel.Secret{}
		errCreate := su.Create(ctx, secret1)
		assert.True(t, errCreate.IsOk())

		// CreateNetwork a second secret
		ctx = ctx.WithValue(share.RequestContextKey, dto.CreateSecretRequest{
			ProjectID:   curProject.ID.Hex(),
			Name:        "test-secret-2",
			Description: "A test secret 2",
			Value:       []byte("test2"),
		})
		secret2 := &secretModel.Secret{}
		errCreate2 := su.Create(ctx, secret2)
		assert.True(t, errCreate2.IsOk())

		// List secrets
		secrets := &[]secretModel.Secret{}
		errList := su.List(ctx, secrets)
		assert.True(t, errList.IsOk())
		assert.Equal(t, 2, len(*secrets))
		for _, s := range *secrets {
			assert.True(t, s.Equals(*secret1) || s.Equals(*secret2))
		}

	})

}

func Test_secretUseCase_Update(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()

	pu := usecase.NewProjectUseCase(repository.NewProjectRepository())
	curProject := model.NewProject("test", "test")
	ctx := NewContext()
	err := pu.Create(ctx, curProject)
	require.True(t, err.IsOk())
	ctx = ctx.WithValue(share.ProjectIDKey, curProject.ID.Hex())

	t.Run("Update existing secret", func(t *testing.T) {
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())

		// CreateNetwork a secret
		ctx = ctx.WithValue(share.RequestContextKey, dto.CreateSecretRequest{
			ProjectID:   curProject.ID.Hex(),
			Name:        "test-secret-1",
			Description: "A test secret 1",
			Value:       []byte("test1"),
		})
		secret := &secretModel.Secret{}
		errCreate := su.Create(ctx, secret)
		assert.True(t, errCreate.IsOk())

		// UpdateNetwork the secret
		ctx = ctx.WithValue(share.RequestContextKey, dto.UpdateSecretRequest{
			Name:        "test-secret-1",
			Description: "A test secret 1 updated",
			Value:       []byte("test1-updated"),
		})
		errUpdate := su.Update(ctx, secret)

		assert.True(t, errUpdate.IsOk())
		assert.Equal(t, "A test secret 1 updated", secret.Description)
		kubeSecret := &corev1.Secret{}
		errKube := kubernetes.Client().Get(ctx, types.NamespacedName{
			Name:      secret.Kubernetes.SecretName,
			Namespace: curProject.Namespace,
		}, kubeSecret)

		assert.Nil(t, errKube)
		assert.Equal(t, "test1-updated", string(kubeSecret.Data[defaultSecretKey]))
	})

	t.Run("Update non-existing secret should return NotFound error", func(t *testing.T) {
		su := usecase.NewSecretUseCase(secretRepo.NewSecretRepository(), secretRepo.NewKubernetesRepository())
		secret := &secretModel.Secret{}
		ctx = ctx.WithValue(share.RequestContextKey, dto.UpdateSecretRequest{
			Name:        "this-secret-does-not-exist",
			Description: "A secret that does not exist",
			Value:       []byte("secret-value"),
		})
		errUpdate := su.Update(ctx, secret)

		assert.Equal(t, errors.NotFound.Code, errUpdate.Code)
	})
}

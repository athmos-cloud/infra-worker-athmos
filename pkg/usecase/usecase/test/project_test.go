package test

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	repository2 "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func Test_projectUseCase_Create(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()

	t.Run("Should successfully create a project and namespace", func(t *testing.T) {
		ctx := NewContext()
		proj := model.NewProject("test-1", "1")
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		err := pu.Create(ctx, proj)

		assert.True(t, err.IsOk())
		assert.False(t, proj.ID.IsZero())
		// Check if namespace has been created
		ns := &corev1.Namespace{}
		errKube := kubernetes.Client().Get(ctx, types.NamespacedName{Name: proj.Namespace}, ns)
		assert.Nil(t, errKube)
	})

	t.Run("Should return Conflict error when create a project if already exists", func(t *testing.T) {
		ctx := NewContext()
		proj := model.NewProject("test-2", "1")
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		err := pu.Create(ctx, proj)
		assert.True(t, err.IsOk())
		err = pu.Create(ctx, proj)
		assert.Equal(t, err.Code, errors.Conflict.Code)
	})
}

func Test_projectUseCase_Delete(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()

	t.Run("Should successfully delete a project and namespace", func(t *testing.T) {
		ctx := NewContext()
		proj := model.NewProject("test-1", "1")
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		err := pu.Create(ctx, proj)
		assert.True(t, err.IsOk())
		ctx = ctx.WithValue(context.ProjectIDKey, proj.ID.Hex())
		err = pu.Delete(ctx, proj)
		assert.Equal(t, err.Code, errors.NoContent.Code)
	})
	t.Run("Should return NotFound error when delete a project if not exists", func(t *testing.T) {
		ctx := NewContext()
		proj := model.NewProject("test-2", "1")
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		ctx = ctx.WithValue(context.ProjectIDKey, primitive.NewObjectID().Hex())
		err := pu.Delete(ctx, proj)
		assert.Equal(t, err.Code, errors.NotFound.Code)
	})
}

func Test_projectUseCase_Get(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()

	t.Run("Should successfully get a project", func(t *testing.T) {
		ctx := NewContext()
		proj := model.NewProject("test-1", "1")
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		err := pu.Create(ctx, proj)
		assert.True(t, err.IsOk())
		ctx = ctx.WithValue(context.ProjectIDKey, proj.ID)
		errGet := pu.Get(ctx, proj)
		assert.True(t, errGet.IsOk())
	})
	t.Run("Should return NotFound error when get a project if not exists", func(t *testing.T) {
		ctx := NewContext()
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		proj := &model.Project{}
		ctx = ctx.WithValue(context.ProjectIDKey, primitive.NewObjectID().Hex())
		errGet := pu.Get(ctx, proj)
		assert.Equal(t, errGet.Code, errors.NotFound.Code)
	})
}

func Test_projectUseCase_List(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()

	t.Run("Should successfully list a project", func(t *testing.T) {
		ctx := NewContext()
		proj1 := model.NewProject("test-1", "1")
		proj2 := model.NewProject("test-2", "1")
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		err1 := pu.Create(ctx, proj1)
		assert.True(t, err1.IsOk())
		err2 := pu.Create(ctx, proj2)
		assert.True(t, err2.IsOk())
		res := &[]model.Project{}
		ctx = ctx.WithValue(context.OwnerIDKey, "1")
		errList := pu.List(ctx, res)
		assert.True(t, errList.IsOk())
		assert.Equal(t, len(*res), 2)
	})

	t.Run("Should return empty list when owner has no project", func(t *testing.T) {
		ctx := NewContext()
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		res := &[]model.Project{}
		ctx = ctx.WithValue(context.OwnerIDKey, "1")
		errList := pu.List(ctx, res)
		assert.True(t, errList.IsOk())
		assert.Equal(t, len(*res), 0)
	})
}

func Test_projectUseCase_Update(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()

	t.Run("Should successfully update an existing project", func(t *testing.T) {
		ctx := NewContext()
		proj := model.NewProject("test-1", "1")
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		err := pu.Create(ctx, proj)
		assert.True(t, err.IsOk())
		ctx = ctx.WithValue(context.ProjectIDKey, proj.ID.Hex())
		ctx = ctx.WithValue(context.RequestKey, dto.UpdateProjectRequest{Name: "test-2"})
		err = pu.Update(ctx, proj)
		assert.True(t, err.IsOk())
		projGet := &model.Project{}
		errGet := pu.Get(ctx, projGet)
		assert.True(t, errGet.IsOk())
		assert.Equal(t, projGet.Name, "test-2")
	})
	t.Run("Should return NotFound error when update a project if not exists", func(t *testing.T) {
		ctx := NewContext()
		proj := model.NewProject("test-1", "1")
		pu := usecase.NewProjectUseCase(repository2.NewProjectRepository())
		ctx = ctx.WithValue(context.ProjectIDKey, primitive.NewObjectID().Hex())
		err := pu.Update(ctx, proj)
		assert.Equal(t, err.Code, errors.NotFound.Code)
	})
}

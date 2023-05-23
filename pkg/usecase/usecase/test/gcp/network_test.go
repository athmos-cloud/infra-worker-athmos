package gcp

import (
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_networkUseCase_Create(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Create a valid network", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a network with a non-existing secret should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a network with an already existing name should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_networkUseCase_Delete(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Delete a valid network should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing network should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a network with children should fail", func(t *testing.T) {

	})
	t.Run("Delete cascade a network should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_networkUseCase_Get(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Get a valid network should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing network should fail", func(t *testing.T) {
		t.Skip("TODO")
	})

}

func Test_networkUseCase_Update(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Update a valid network should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Update a non-existing network should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

package gcp

import (
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_subnetworkUseCase_Create(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Create a valid subnetwork", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a subnetwork with a non-existing secret should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a subnetwork with an already existing name should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_subnetworkUseCase_Delete(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Delete a valid subnetwork should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing subnetwork should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a subnetwork with children should fail", func(t *testing.T) {

	})
	t.Run("Delete cascade a subnetwork should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_subnetworkUseCase_Get(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Get a valid subnetwork should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing subnetwork should fail", func(t *testing.T) {
		t.Skip("TODO")
	})

}

func Test_subnetworkUseCase_Update(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Update a valid subnetwork should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Update a non-existing subnetwork should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

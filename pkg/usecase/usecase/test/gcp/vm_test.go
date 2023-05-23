package gcp

import (
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_vmUseCase_Create(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Create a valid vm", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a vm with a non-existing secret should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a vm with an already existing name should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_vmUseCase_Delete(t *testing.T) {
	mongoC, kubeC := Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		require.NoError(t, gnomock.Stop(kubeC))
	}()
	t.Run("Delete a valid vm should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing vm should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a vm with children should fail", func(t *testing.T) {

	})
	t.Run("Delete cascade a vm should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

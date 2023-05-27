package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_firewallUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("Create a valid firewall", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a firewall with a non-existing secret should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a firewall with an already existing name should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_firewallUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("Delete a valid firewall should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing firewall should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a firewall with children should fail", func(t *testing.T) {

	})
	t.Run("Delete cascade a firewall should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_firewallUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("Get a valid firewall should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing firewall should fail", func(t *testing.T) {
		t.Skip("TODO")
	})

}

func Test_firewallUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("Update a valid firewall should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Update a non-existing firewall should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

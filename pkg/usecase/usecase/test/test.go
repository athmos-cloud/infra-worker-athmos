package test

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/kamva/mgm/v3"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/k3s"
	"github.com/orlangure/gnomock/preset/mongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func Init(t *testing.T) (*gnomock.Container, *gnomock.Container) {
	mongoC := InitMongo(t)
	kubeC := InitKubernetes(t)

	return mongoC, kubeC
}

func InitMongo(t *testing.T) *gnomock.Container {
	user := "test"
	password := "test"
	db := "infra-test"
	p := mongo.Preset(
		mongo.WithData("./testdata/"),
		mongo.WithUser(user, password),
	)
	c, err := gnomock.Start(p)
	require.NoError(t, err)
	addr := c.DefaultAddress()
	uri := fmt.Sprintf("mongodb://%s:%s@%s", user, password, addr)

	errMgm := mgm.SetDefaultConfig(nil, db, options.Client().ApplyURI(uri))
	require.NoError(t, errMgm)

	return c
}

func InitKubernetes(t *testing.T) *gnomock.Container {
	t.Parallel()
	p := k3s.Preset()
	c, err := gnomock.Start(
		p,
		gnomock.WithContainerName("k3s"),
	)
	require.NoError(t, err)
	kubeConfig, err := k3s.Config(c)
	require.NoError(t, err)
	k8sCli, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)
	kubernetes.SetClient(k8sCli)
	require.NoError(t, err)

	return c
}

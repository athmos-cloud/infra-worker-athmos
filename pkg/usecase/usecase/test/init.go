package test

import (
	"context"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/kamva/mgm/v3"
	helmclient "github.com/mittwald/go-helm-client"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/k3s"
	"github.com/orlangure/gnomock/preset/mongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/options"
	"helm.sh/helm/v3/pkg/repo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

const (
	crossplaneNamespace        = "crossplane"
	crossplaneChartRepoAddress = "https://charts.crossplane.io/stable"
	crossplaneRepoName         = "crossplane-stable"
	crossplaneChartName        = "crossplane"
	crossplaneReleaseName      = "crossplane"
	CleanupOnFail              = true
	RepositoryCache            = "/tmp/.helmcache"
	RepositoryConfig           = "/tmp/.helmconfig"
	DefaultNamespace           = "default"
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
func InitCrossplane(t *testing.T, kubeConfig []byte) {
	t.Parallel()
	ctx := context.Background()
	err := kubernetes.Client().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: crossplaneNamespace,
		},
	})
	require.NoError(t, err)
	cli, err := helmclient.NewClientFromKubeConf(
		&helmclient.KubeConfClientOptions{
			KubeContext: "",
			KubeConfig:  kubeConfig,
			Options: &helmclient.Options{
				RepositoryCache:  RepositoryCache,
				RepositoryConfig: RepositoryConfig,
				Namespace:        DefaultNamespace,
			},
		})
	require.NoError(t, err)
	err = cli.AddOrUpdateChartRepo(
		repo.Entry{
			Name: crossplaneRepoName,
			URL:  crossplaneChartRepoAddress,
		})
	require.NoError(t, err)
	_, err = cli.InstallChart(
		ctx,
		&helmclient.ChartSpec{
			ChartName:     fmt.Sprintf("%s/%s", crossplaneRepoName, crossplaneChartName),
			ReleaseName:   crossplaneReleaseName,
			CleanupOnFail: CleanupOnFail,
		},
		&helmclient.GenericHelmOptions{},
	)
	require.NoError(t, err)
}

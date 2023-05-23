package gcp

import (
	"context"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	"github.com/crossplane/crossplane/apis/pkg/meta/v1"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/k3s"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

const (
	gcpCrossplaneProviderName   = "provider-gcp"
	gcpCrossplanePackage        = "xpkg.upbound.io/upbound/provider-gcp"
	gcpCrossplanePackageVersion = "v0.32.0"
)

func Init(t *testing.T) (*gnomock.Container, *gnomock.Container) {
	mongoC, kubeC := test.Init(t)

	kubeConfig, err := k3s.Config(kubeC)
	require.NoError(t, err)
	test.InitCrossplane(t, []byte(kubeConfig.String()))
	providerImage := fmt.Sprintf("%s:%s", gcpCrossplanePackage, gcpCrossplanePackageVersion)
	provider := v1.Provider{
		ObjectMeta: metav1.ObjectMeta{
			Name: gcpCrossplaneProviderName,
		},
		Spec: v1.ProviderSpec{
			Controller: v1.ControllerSpec{
				Image: &providerImage,
			},
		},
	}
	err = kubernetes.Client().Create(context.Background(), &provider)
	require.NoError(t, err)

	return mongoC, kubeC
}

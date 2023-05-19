package gcp

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	gcpRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/gcp"
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	autoCreateSubnetworks = false
)

type networkRepository struct{}

func NewNetworkRepository() gcpRepo.Network {
	return &networkRepository{}
}

func (nr *networkRepository) Find(ctx context.Context, opt option.Option) *resource.Network {
	//TODO implement me
	panic("implement me")
}

func (nr *networkRepository) FindAll(ctx context.Context, opt option.Option) []*resource.Network {
	//TODO implement me
	panic("implement me")
}

func (nr *networkRepository) Create(ctx context.Context, network *resource.Network) *resource.Network {
	gcpNetwork := v1beta1.Network{
		ObjectMeta: metav1.ObjectMeta{
			Name: network.Identifier.NetworkID,
		},
		Spec: v1beta1.NetworkSpec{
			ForProvider: v1beta1.NetworkParameters{
				Project:               &network.Identifier.VPCID,
				AutoCreateSubnetworks: &autoCreateSubnetworks,
			},
		},
	}
	if err := kubernetes.Client().K8sClient.Create(ctx, &gcpNetwork); err != nil {
		panic(errors.KubernetesError.WithMessage(err.Error()))
	}
	return network
}

func (nr *networkRepository) Update(ctx context.Context, network *resource.Network) *resource.Network {
	//TODO implement me
	panic("implement me")
}

func (nr *networkRepository) Delete(ctx context.Context, network *resource.Network) {
	//TODO implement me
	panic("implement me")
}

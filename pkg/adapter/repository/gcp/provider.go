package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	domainTypes "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/samber/lo"
	"github.com/upbound/provider-gcp/apis/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (gcp *gcpRepository) FindProvider(ctx context.Context, opt option.Option) (*resource.Provider, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpProvider := &v1beta1.ProviderConfig{}
	if err := kubernetes.Client().Get(ctx, types.NamespacedName{Name: req.Name}, gcpProvider); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get provider %s", req.Name))
	}
	mod, err := gcp.toModelProvider(gcpProvider)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (gcp *gcpRepository) FindAllProviders(ctx context.Context, opt option.Option) (*resource.ProviderCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpProviders := &v1beta1.ProviderConfigList{}
	if err := kubernetes.Client().List(ctx, gcpProviders, &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("providers not found"))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list providers"))
	}
	mod, err := gcp.toModelProviderCollection(gcpProviders)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (gcp *gcpRepository) FindProviderStack(ctx context.Context, opt option.Option) (*resource.Provider, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpProviders := &v1beta1.ProviderConfig{}
	if err := kubernetes.Client().Get(ctx, types.NamespacedName{Name: req.Name}, gcpProviders); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get provider %s", req.Name))
	}
	providerModel, err := gcp.toModelProvider(gcpProviders)
	if !err.IsOk() {
		return nil, err
	}
	// Find Networks recursively
	networks, err := gcp.FindAllRecursiveNetworks(ctx, option.Option{
		Value: resourceRepo.FindAllResourceOption{
			Labels: map[string]string{
				identifier.ProviderLabelKey: providerModel.IdentifierID.Provider,
			},
			Namespace: req.Namespace,
		},
	}, nil)
	if !err.IsOk() {
		return nil, err
	}
	providerModel.Networks = *networks

	return nil, errors.OK
}

func (gcp *gcpRepository) CreateProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	gcpProvider := gcp.toGCPProvider(ctx, provider)
	if err := kubernetes.Client().Create(ctx, gcpProvider); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("provider %s already exists", provider.IdentifierID.Provider))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s", provider.IdentifierID.Provider))
	}

	return errors.Created
}

func (gcp *gcpRepository) UpdateProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	gcpProvider := gcp.toGCPProvider(ctx, provider)
	if err := kubernetes.Client().Update(ctx, gcpProvider); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", provider.IdentifierID.Provider))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s", provider.IdentifierID.Provider))
	}

	return errors.NoContent
}

func (gcp *gcpRepository) DeleteProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	gcpSubnetwork := gcp.toGCPProvider(ctx, provider)
	if err := kubernetes.Client().Delete(ctx, gcpSubnetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", provider.IdentifierID.Provider))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s", provider.IdentifierID.Provider))
	}

	return errors.NoContent
}

func (gcp *gcpRepository) DeleteProviderCascade(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) toModelProvider(provider *v1beta1.ProviderConfig) (*resource.Provider, errors.Error) {
	id := identifier.Provider{}
	if err := id.FromLabels(provider.Labels); !err.IsOk() {
		return nil, err
	}
	return &resource.Provider{
		IdentifierID: id,
		IdentifierName: identifier.Provider{
			Provider: provider.ObjectMeta.Labels[identifier.ProviderLabelKey],
			VPC:      provider.Spec.ProjectID,
		},
		Type: domainTypes.ProviderGCP,
		Auth: secret.Secret{
			Name:        provider.ObjectMeta.Labels[secret.NameLabelKey],
			Description: provider.ObjectMeta.Labels[secret.DescriptionLabelKey],
			Kubernetes: secret.Kubernetes{
				SecretName: provider.Spec.Credentials.CommonCredentialSelectors.SecretRef.SecretReference.Name,
				Namespace:  provider.Spec.Credentials.CommonCredentialSelectors.SecretRef.SecretReference.Namespace,
				SecretKey:  provider.Spec.Credentials.CommonCredentialSelectors.SecretRef.Key,
			},
		},
	}, errors.OK
}

func (gcp *gcpRepository) toGCPProvider(ctx context.Context, provider *resource.Provider) *v1beta1.ProviderConfig {
	resLabels := lo.Assign(
		crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)),
		provider.IdentifierID.ToLabels(),
		provider.IdentifierName.GetLabelName(),
		map[string]string{
			secret.NameLabelKey:        provider.Auth.Name,
			secret.DescriptionLabelKey: provider.Auth.Description,
		},
	)
	return &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      provider.IdentifierID.Provider,
			Namespace: provider.IdentifierID.Provider,
			Labels:    resLabels,
		},
		Spec: v1beta1.ProviderConfigSpec{
			ProjectID: provider.IdentifierID.VPC,
			Credentials: v1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: provider.Auth.Kubernetes.SecretKey,
						SecretReference: xpv1.SecretReference{
							Namespace: provider.Auth.Kubernetes.Namespace,
							Name:      provider.Auth.Kubernetes.SecretName,
						},
					},
				},
			},
		},
	}
}

func (gcp *gcpRepository) toModelProviderCollection(providers *v1beta1.ProviderConfigList) (*resource.ProviderCollection, errors.Error) {
	providerCollection := resource.ProviderCollection{}
	for _, provider := range providers.Items {
		modelProvider, err := gcp.toModelProvider(&provider)
		if !err.IsOk() {
			return nil, err
		}
		providerCollection[modelProvider.IdentifierName.Provider] = *modelProvider
	}
	return &providerCollection, errors.OK
}

package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	modelTypes "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
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
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name}, gcpProvider); err != nil {
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
	if err := kubernetes.Client().Client.List(ctx, gcpProviders, &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}); err != nil {
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
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name}, gcpProviders); err != nil {
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
			Labels: providerModel.IdentifierID.ToIDLabels(),
		},
	}, nil)
	if !err.IsOk() {
		return nil, err
	}
	providerModel.Networks = *networks

	return providerModel, errors.OK
}

func (gcp *gcpRepository) CreateProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	if exists, err := gcp.ProviderExists(ctx, provider.IdentifierName); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("provider %s already exists", provider.IdentifierName.Provider))
	}
	gcpProvider := gcp.toGCPProvider(ctx, provider)
	if err := kubernetes.Client().Client.Create(ctx, gcpProvider); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("provider %s already exists", provider.IdentifierID.Provider))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s", provider.IdentifierID.Provider))
	}

	return errors.Created
}

func (gcp *gcpRepository) UpdateProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	existingProvider := &v1beta1.ProviderConfig{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: provider.IdentifierID.Provider}, existingProvider); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", provider.IdentifierID.Provider))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get provider %s", provider.IdentifierID.Provider))
	}

	gcpProvider := gcp.toGCPProvider(ctx, provider)
	existingProvider.Spec = gcpProvider.Spec
	existingProvider.Labels = gcpProvider.Labels

	if err := kubernetes.Client().Client.Update(ctx, existingProvider); err != nil {
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s", provider.IdentifierID.Provider))
	}

	return errors.NoContent
}

func (gcp *gcpRepository) DeleteProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	gcpSubnetwork := gcp.toGCPProvider(ctx, provider)

	if err := kubernetes.Client().Client.Delete(ctx, gcpSubnetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", provider.IdentifierID.Provider))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s", provider.IdentifierID.Provider))
	}
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, provider.IdentifierID.ToIDLabels())
	networks, networksErr := gcp.FindAllNetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	if !networksErr.IsOk() {
		return networksErr
	}
	if len(*networks) > 0 {
		return errors.Conflict.WithMessage(fmt.Sprintf("provider %s still has networks", provider.IdentifierID.Provider))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteProviderCascade(ctx context.Context, provider *resource.Provider) errors.Error {
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, provider.IdentifierID.ToIDLabels())
	networks, networksErr := gcp.FindAllNetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	if !networksErr.IsOk() {
		return networksErr
	}

	for _, network := range *networks {
		if networkErr := gcp.DeleteNetworkCascade(ctx, &network); !networkErr.IsOk() {
			return networkErr
		}
	}

	return gcp.DeleteProvider(ctx, provider)
}

func (gcp *gcpRepository) ProviderExists(ctx context.Context, name identifier.Provider) (bool, errors.Error) {
	providerLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, name.ToNameLabels())
	gcpProviders := &v1beta1.ProviderConfigList{}

	if err := kubernetes.Client().Client.List(ctx, gcpProviders, client.MatchingLabels(providerLabels)); err != nil {
		return false, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list providers"))
	}
	return len(gcpProviders.Items) > 0, errors.OK
}

func (gcp *gcpRepository) toModelProvider(provider *v1beta1.ProviderConfig) (*resource.Provider, errors.Error) {
	id := identifier.Provider{}
	name := identifier.Provider{}
	if err := id.IDFromLabels(provider.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(provider.Labels); !err.IsOk() {
		return nil, err
	}
	return &resource.Provider{
		IdentifierID:   id,
		IdentifierName: name,
		Type:           modelTypes.ProviderGCP,
		Auth: resource.ProviderAuth{
			Name: provider.ObjectMeta.Labels[secret.NameLabelKey],
			KubernetesSecret: secret.Kubernetes{
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
		provider.IdentifierID.ToIDLabels(),
		provider.IdentifierName.ToNameLabels(),
		map[string]string{
			secret.NameLabelKey: provider.Auth.Name,
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
						Key: provider.Auth.KubernetesSecret.SecretKey,
						SecretReference: xpv1.SecretReference{
							Namespace: provider.Auth.KubernetesSecret.Namespace,
							Name:      provider.Auth.KubernetesSecret.SecretName,
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

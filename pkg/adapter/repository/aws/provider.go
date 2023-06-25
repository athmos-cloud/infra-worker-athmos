package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	modelTypes "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/samber/lo"
	"github.com/upbound/provider-aws/apis/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (aws *awsRepository) FindProvider(ctx context.Context, opt option.Option) (*resource.Provider, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	awsProvider := &v1beta1.ProviderConfig{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name}, awsProvider); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get provider %s", req.Name))
	}
	mod, err := aws.toModelProvider(awsProvider)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (aws *awsRepository) FindAllProviders(ctx context.Context, opt option.Option) (*resource.ProviderCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	awsProviders := &v1beta1.ProviderConfigList{}
	if err := kubernetes.Client().Client.List(ctx, awsProviders, &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list providers"))
	}
	mod, err := aws.toModelProviderCollection(awsProviders)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (aws *awsRepository) FindProviderStack(ctx context.Context, opt option.Option) (*resource.Provider, errors.Error) {
	//TODO
	panic("Implement me")
}

func (aws *awsRepository) CreateProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	if exists, err := aws.ProviderExists(ctx, provider.IdentifierName); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("provider %s already exists", provider.IdentifierName.Provider))
	}
	awsProvider := aws.toAWSProvider(ctx, provider)
	if err := kubernetes.Client().Client.Create(ctx, awsProvider); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("provider %s already exists", provider.IdentifierID.Provider))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s", provider.IdentifierID.Provider))
	}

	return errors.Created
}

func (aws *awsRepository) UpdateProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	existingProvider := &v1beta1.ProviderConfig{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: provider.IdentifierID.Provider}, existingProvider); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", provider.IdentifierID.Provider))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get provider %s", provider.IdentifierID.Provider))
	}

	awsProvider := aws.toAWSProvider(ctx, provider)
	existingProvider.Spec = awsProvider.Spec
	existingProvider.Labels = awsProvider.Labels

	if err := kubernetes.Client().Client.Update(ctx, existingProvider); err != nil {
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s", provider.IdentifierID.Provider))
	}

	return errors.NoContent
}

func (aws *awsRepository) DeleteProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	awsProvider := aws.toAWSProvider(ctx, provider)
	if err := kubernetes.Client().Client.Delete(ctx, awsProvider); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", provider.IdentifierID.Provider))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s", provider.IdentifierID.Provider))
	}

	return errors.NoContent
}

func (aws *awsRepository) DeleteProviderCascade(ctx context.Context, provider *resource.Provider) errors.Error {
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, provider.IdentifierID.ToIDLabels())
	networks, err := aws.FindAllNetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	if !err.IsOk() {
		return err
	}

	for _, network := range *networks {
		if err := aws.DeleteNetworkCascade(ctx, &network); !err.IsOk() {
			return err
		}
	}

	return aws.DeleteProvider(ctx, provider)
}

func (aws *awsRepository) ProviderExists(ctx context.Context, name identifier.Provider) (bool, errors.Error) {
	providerLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, name.ToNameLabels())
	awsProviders := &v1beta1.ProviderConfigList{}

	if err := kubernetes.Client().Client.List(ctx, awsProviders, client.MatchingLabels(providerLabels)); err != nil {
		return false, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list providers"))
	}
	return len(awsProviders.Items) > 0, errors.OK
}

func (aws *awsRepository) toModelProvider(provider *v1beta1.ProviderConfig) (*resource.Provider, errors.Error) {
	id := identifier.Provider{}
	name := identifier.Provider{}
	if err := id.IDFromLabels(provider.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(provider.Labels); !err.IsOk() {
		return nil, err
	}
	return &resource.Provider{
		Metadata: metadata.Metadata{
			Status: metadata.StatusFromKubernetesStatus(provider.Status.Conditions),
		},
		IdentifierID:   id,
		IdentifierName: name,
		Type:           modelTypes.ProviderAWS,
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

func (aws *awsRepository) toAWSProvider(ctx context.Context, provider *resource.Provider) *v1beta1.ProviderConfig {
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

func (aws *awsRepository) toModelProviderCollection(providers *v1beta1.ProviderConfigList) (*resource.ProviderCollection, errors.Error) {
	providerCollection := resource.ProviderCollection{}
	for _, provider := range providers.Items {
		modelProvider, err := aws.toModelProvider(&provider)
		if !err.IsOk() {
			return nil, err
		}
		providerCollection[modelProvider.IdentifierName.Provider] = *modelProvider
	}
	return &providerCollection, errors.OK
}

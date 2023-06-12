package aws

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/samber/lo"
	"github.com/upbound/provider-aws/apis/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (aws *awsRepository) FindProvider(ctx context.Context, opt option.Option) (*resource.Provider, errors.Error) {
	//TODO
	panic("Implement me")
}

func (aws *awsRepository) FindAllProviders(ctx context.Context, opt option.Option) (*resource.ProviderCollection, errors.Error) {
	//TODO
	panic("Implement me")
}

func (aws *awsRepository) FindProviderStack(ctx context.Context, opt option.Option) (*resource.Provider, errors.Error) {
	//TODO
	panic("Implement me")
}

func (aws *awsRepository) CreateProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO
	panic("Implement me")
}

func (aws *awsRepository) UpdateProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO
	panic("Implement me")
}

func (aws *awsRepository) DeleteProvider(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO
	panic("Implement me")
}

func (aws *awsRepository) DeleteProviderCascade(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO
	panic("Implement me")
}

func (aws *awsRepository) ProviderExists(ctx context.Context, name identifier.Provider) (bool, errors.Error) {
	//TODO
	panic("Implement me")
}

func (aws *awsRepository) toModelProvider(provider *v1beta1.ProviderConfig) (*resource.Provider, errors.Error) {
	//TODO
	panic("Implement me")
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

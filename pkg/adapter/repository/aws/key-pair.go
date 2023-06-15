package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/samber/lo"
	"github.com/upbound/provider-aws/apis/ec2/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (aws *awsRepository) _createKeyPair(ctx context.Context, vm *instance.VM) errors.Error {
	awsKeyPair, keyErr := aws._toAwsKeyPair(ctx, vm)
	if !keyErr.IsOk() {
		return keyErr
	}

	if err := kubernetes.Client().Client.Create(ctx, awsKeyPair); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("SSH key pair for vm %s already exists", vm.IdentifierName.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create SSH key pair for vm %s", vm.IdentifierName.VM))
	}
	return errors.Created
}

func (aws *awsRepository) _getKeyPair(ctx context.Context, vm *v1beta1.Instance) (*v1beta1.KeyPair, errors.Error) {
	awsKeyPair := &v1beta1.KeyPair{}
	name := fmt.Sprintf("%s-keypair", vm.Name)
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: name}, awsKeyPair); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf(
				"key pair %s not found in namespace %s",
				name,
				vm.Namespace))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf(
			"unable to get key pair %s in namespace %s",
			name,
			vm.Namespace))
	}

	return awsKeyPair, errors.OK
}

func (aws *awsRepository) _updateKeyPair(ctx context.Context, vm *instance.VM) errors.Error {
	awsKeyPair, keyErr := aws._toAwsKeyPair(ctx, vm)
	if !keyErr.IsOk() {
		return keyErr
	}

	if err := kubernetes.Client().Client.Update(ctx, awsKeyPair); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("key pair %s not found", awsKeyPair.Name))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update key pair %s", awsKeyPair.Name))
	}
	return errors.NoContent
}

func (aws *awsRepository) _deleteKeyPair(ctx context.Context, vm *instance.VM) errors.Error {
	awsKeyPair, keyPairErr := aws._toAwsKeyPair(ctx, vm)
	if !keyPairErr.IsOk() {
		return keyPairErr
	}

	if err := kubernetes.Client().Client.Delete(ctx, awsKeyPair); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(
				fmt.Sprintf(
					"keypair %s not found in namespace %s",
					awsKeyPair.Name,
					awsKeyPair.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf(
			"unable to delete vm %s in namespace %s",
			awsKeyPair.Name,
			awsKeyPair.Namespace))
	}
	return errors.NoContent
}

func (aws *awsRepository) _toAwsKeyPair(ctx context.Context, vm *instance.VM) (*v1beta1.KeyPair, errors.Error) {
	if len(vm.Auths) == 0 {
		return &v1beta1.KeyPair{}, errors.BadRequest.WithMessage("Missing SSH key pair")
	}
	modelKeyPair := vm.Auths[0]

	sshKeysLabels := crossplane.ToSSHKeySecretLabels(vm.Auths)
	keyPairLabels := lo.Assign(
		crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)),
		vm.IdentifierID.ToIDLabels(),
		vm.IdentifierName.ToNameLabels(),
		sshKeysLabels)

	return &v1beta1.KeyPair{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-keypair", vm.IdentifierID.VM),
			Labels:      keyPairLabels,
			Annotations: crossplane.GetAnnotations(vm.Metadata.Managed, vm.IdentifierName.Network),
		},
		Spec: v1beta1.KeyPairSpec{
			ForProvider: v1beta1.KeyPairParameters{
				PublicKey: &modelKeyPair.PublicKey,
				Region:    &vm.Zone,
			},
		},
	}, errors.OK
}

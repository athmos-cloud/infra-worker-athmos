package secret

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	"github.com/kamva/mgm/v3"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
)

type kubernetesRepository struct{}

func NewKubernetesRepository() repository.KubernetesSecret {
	return &kubernetesRepository{}
}

type CreateKubernetesSecretRequest struct {
	ProjectID   string
	SecretName  string
	SecretKey   string
	SecretValue []byte
}

type UpdateKubernetesSecretRequest struct {
	ProjectID   string
	SecretName  string
	SecretKey   string
	SecretValue []byte
}

type DeleteKubernetesSecretRequest struct {
	ProjectID  string
	SecretName string
}

func (k *kubernetesRepository) Create(ctx context.Context, opt option.Option) (*secret.Kubernetes, errors.Error) {
	if !opt.SetType(reflect.TypeOf(CreateKubernetesSecretRequest{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(
			fmt.Sprintf("expected CreateKubernetesSecretRequest option, got %s", reflect.TypeOf(opt.Value).String()),
		)
	}
	req := opt.Value.(CreateKubernetesSecretRequest)
	ns, err := k.getProjectNamespace(req.ProjectID)
	if !err.IsOk() {
		return nil, err
	}
	kubeSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.SecretName,
			Namespace: ns,
		},
		Data: map[string][]byte{
			req.SecretKey: req.SecretValue,
		},
	}
	if errKube := kubernetes.Client().Create(ctx, kubeSecret); errKube != nil {
		return nil, errors.KubernetesError.WithMessage(errKube.Error())
	}

	createdSecret := secret.NewKubernetesSecret(req.SecretName, req.SecretKey, ns)
	return &createdSecret, errors.Created
}

func (k *kubernetesRepository) Update(ctx context.Context, opt option.Option) errors.Error {
	if !opt.SetType(reflect.TypeOf(UpdateKubernetesSecretRequest{}).String()).Validate() {
		return errors.InvalidOption.WithMessage(
			fmt.Sprintf("expected UpdateKubernetesSecretRequest option, got %s", reflect.TypeOf(opt.Value).String()),
		)
	}
	req := opt.Value.(UpdateKubernetesSecretRequest)
	ns, err := k.getProjectNamespace(req.ProjectID)
	if !err.IsOk() {
		return err
	}
	kubeSecret := &corev1.Secret{}
	if errKube := kubernetes.Client().Get(ctx,
		types.NamespacedName{Namespace: ns, Name: req.SecretName},
		kubeSecret); errKube != nil {
		if k8serrors.IsNotFound(errKube) {
			return errors.NotFound.WithMessage(fmt.Sprintf("Secret with name %s not found in namespace %s", req.SecretName, ns))
		}
		return errors.KubernetesError.WithMessage(errKube.Error())
	}
	kubeSecret.Data[req.SecretKey] = req.SecretValue

	if errKube := kubernetes.Client().Update(ctx, kubeSecret); errKube != nil {
		return errors.KubernetesError.WithMessage(errKube.Error())
	}
	return errors.NoContent
}

func (k *kubernetesRepository) Delete(ctx context.Context, opt option.Option) errors.Error {
	if !opt.SetType(reflect.TypeOf(DeleteKubernetesSecretRequest{}).String()).Validate() {
		return errors.InvalidOption.WithMessage(
			fmt.Sprintf("expected DeleteKubernetesSecretRequest option, got %s", reflect.TypeOf(opt.Value).String()),
		)
	}
	req := opt.Value.(DeleteKubernetesSecretRequest)
	ns, err := k.getProjectNamespace(req.ProjectID)
	if !err.IsOk() {
		return err
	}
	kubeSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.SecretName,
			Namespace: ns,
		},
	}
	if errKube := kubernetes.Client().Delete(ctx, kubeSecret); errKube != nil {
		if k8serrors.IsNotFound(errKube) {
			return errors.NotFound.WithMessage(fmt.Sprintf("Secret with name %s not found in namespace %s", req.SecretName, ns))
		}
		return errors.KubernetesError.WithMessage(errKube.Error())
	}
	return errors.NoContent
}

func (k *kubernetesRepository) getProjectNamespace(projectID string) (string, errors.Error) {
	project := &model.Project{}
	err := mgm.Coll(project).FindByID(projectID, project)
	if err != nil {
		return "", errors.NotFound.WithMessage(fmt.Sprintf("Project with id %s not found", projectID))
	}
	return project.Namespace, errors.OK
}

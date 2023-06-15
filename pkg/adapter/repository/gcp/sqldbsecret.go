package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (gcp *gcpRepository) _createSqlPasswordSecret(ctx context.Context, db *instance.SqlDB) errors.Error {
	ns, ok := ctx.Value(context.CurrentNamespace).(string)
	if !ok {
		return errors.InternalError.WithMessage("unable to get current namespace")
	}
	labels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), db.IdentifierID.ToIDLabels(), db.IdentifierName.ToNameLabels())
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      db.IdentifierID.SqlDB,
			Namespace: ns,
			Labels:    labels,
		},
		Data: map[string][]byte{
			"password": []byte(db.Auth.RootPassword),
		},
	}
	if err := kubernetes.Client().Client.Create(ctx, &secret); err != nil {
		return errors.KubernetesError.WithMessage(err.Error())
	}

	return errors.Created
}

func (gcp *gcpRepository) _updateSqlPasswordSecret(ctx context.Context, db *instance.SqlDB) errors.Error {
	ns, ok := ctx.Value(context.CurrentNamespace).(string)
	if !ok {
		return errors.InternalError.WithMessage("unable to get current namespace")
	}
	existingSecret := corev1.Secret{}
	if err := kubernetes.Client().Client.Get(ctx, client.ObjectKey{Name: db.IdentifierID.SqlDB, Namespace: ns}, &existingSecret); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("secret %s not found in namespace %s", db.IdentifierID.SqlDB, ns))
		}
		return errors.KubernetesError.WithMessage(err.Error())
	}
	if string(existingSecret.Data["password"]) == db.Auth.RootPassword {
		return errors.OK
	}
	existingSecret.Data["password"] = []byte(db.Auth.RootPassword)

	if err := kubernetes.Client().Client.Update(ctx, &existingSecret); err != nil {
		return errors.KubernetesError.WithMessage(err.Error())
	}

	return errors.NoContent
}

func (gcp *gcpRepository) _getSqlPasswordSecret(ctx context.Context, db *instance.SqlDB) errors.Error {
	ns, ok := ctx.Value(context.CurrentNamespace).(string)
	if !ok {
		return errors.InternalError.WithMessage("unable to get current namespace")
	}
	existingSecret := corev1.Secret{}
	if err := kubernetes.Client().Client.Get(ctx, client.ObjectKey{Name: db.IdentifierID.SqlDB, Namespace: ns}, &existingSecret); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("secret %s not found in namespace %s", db.IdentifierID.SqlDB, ns))
		}
		return errors.KubernetesError.WithMessage(err.Error())
	}
	db.Auth.RootPassword = string(existingSecret.Data["password"])

	return errors.OK
}

func (gcp *gcpRepository) _deleteSqlPasswordSecret(ctx context.Context, db *instance.SqlDB) errors.Error {
	ns, ok := ctx.Value(context.CurrentNamespace).(string)
	if !ok {
		return errors.InternalError.WithMessage("unable to get current namespace")
	}
	existingSecret := corev1.Secret{}
	if err := kubernetes.Client().Client.Get(ctx, client.ObjectKey{Name: db.IdentifierID.SqlDB, Namespace: ns}, &existingSecret); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("secret %s not found in namespace %s", db.IdentifierID.SqlDB, ns))
		}
		return errors.KubernetesError.WithMessage(err.Error())
	}
	if err := kubernetes.Client().Client.Delete(ctx, &existingSecret); err != nil {
		return errors.KubernetesError.WithMessage(err.Error())
	}
	return errors.NoContent
}

package repository

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	"golang.org/x/crypto/ssh"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

type sshKey struct{}

const (
	SecretPrivateKeyKey = "privateKey"
	SecretPublicKeyKey  = "publicKey"
	SecretUsernameKey   = "username"
	reservedCharacter   = "_"
)

func NewSSHKeyRepository() resourceRepo.SSHKeys {
	return &sshKey{}
}

func (s *sshKey) Create(ctx context.Context, key *model.SSHKey) errors.Error {
	if key.KeyLength == 0 {
		key.KeyLength = 2048
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, key.KeyLength)
	if err != nil {
		return errors.InternalError.WithMessage(err.Error())
	}
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return errors.InternalError.WithMessage(err.Error())
	}
	privateKeyPem := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	privateKeyBytes := pem.EncodeToMemory(privateKeyPem)
	key.PublicKey = string(ssh.MarshalAuthorizedKey(publicKey))
	key.PrivateKey = string(privateKeyBytes)
	secretNameFormatted := strings.ReplaceAll(key.SecretName, reservedCharacter, "-")
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretNameFormatted,
			Namespace: key.SecretNamespace,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			SecretUsernameKey:   []byte(key.Username),
			SecretPrivateKeyKey: privateKeyBytes,
			SecretPublicKeyKey:  ssh.MarshalAuthorizedKey(publicKey),
		},
	}
	if errKube := kubernetes.Client().Client.Create(ctx, secret); errKube != nil {
		return errors.KubernetesError.WithMessage(errKube.Error())
	}

	return errors.Created
}

func (s *sshKey) CreateList(ctx context.Context, list model.SSHKeyList) errors.Error {
	for _, key := range list {
		if err := s.Create(ctx, key); !err.IsOk() {
			return err
		}
	}
	return errors.Created
}

func (s *sshKey) Get(ctx context.Context, key *model.SSHKey) errors.Error {
	secret := &corev1.Secret{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Namespace: key.SecretNamespace, Name: key.SecretName}, secret); err != nil {
		return errors.KubernetesError.WithMessage(err.Error())
	}
	key.Username = string(secret.Data[SecretUsernameKey])
	key.PublicKey = string(secret.Data[SecretPublicKeyKey])
	key.PrivateKey = string(secret.Data[SecretPrivateKeyKey])
	return errors.OK

}

func (s *sshKey) GetList(c context.Context, list model.SSHKeyList) errors.Error {
	for _, key := range list {
		if err := s.Get(c, key); !err.IsOk() {
			return err
		}
	}
	return errors.OK
}

func (s *sshKey) Delete(ctx context.Context, key *model.SSHKey) errors.Error {
	secret := &corev1.Secret{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Namespace: key.SecretNamespace, Name: key.SecretName}, secret); err != nil {
		return errors.KubernetesError.WithMessage(err.Error())
	}
	if err := kubernetes.Client().Client.Delete(ctx, secret); err != nil {
		return errors.KubernetesError.WithMessage(err.Error())
	}
	return errors.OK
}

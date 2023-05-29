package repository

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	"golang.org/x/crypto/ssh"
)

type sshKey struct{}

func NewSSHKey() resourceRepo.SSHKeys {
	return &sshKey{}
}

func (s *sshKey) Create(ctx context.Context, key *model.SSHKey) errors.Error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
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

	return errors.Created
}

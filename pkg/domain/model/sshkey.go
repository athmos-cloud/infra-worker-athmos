package model

type SSHKey struct {
	KeyLength       int
	PublicKey       string
	PrivateKey      string
	Username        string
	SecretName      string
	SecretNamespace string
}

type SSHKeyList []*SSHKey

package crossplane

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"strings"
)

const (
	managedByLabel           = "app.kubernetes.io/managed-by"
	managedByValue           = "athmos"
	VMSSHKeysSecretNamespace = "vm-ssh-keys-secret-namespace"
	VMSSHKeysNames           = "vm-ssh-keys-names"
	VMSSHKeysNamesSeparator  = "."
	VMPublicIPLabel          = "vm-has-public-ip"
)

func GetBaseLabels(projectID string) map[string]string {
	return map[string]string{
		managedByLabel:          managedByValue,
		model.ProjectIDLabelKey: projectID,
	}
}

func ToSSHKeySecretLabels(keyList model.SSHKeyList) map[string]string {
	if len(keyList) == 0 {
		return map[string]string{}
	}
	sshKeyLabels := map[string]string{
		VMSSHKeysSecretNamespace: keyList[0].SecretNamespace,
	}
	for _, key := range keyList {
		sshKeyLabels[VMSSHKeysNames] = sshKeyLabels[VMSSHKeysNames] + key.SecretName + VMSSHKeysNamesSeparator
	}
	sshKeyLabels[VMSSHKeysNames] = strings.TrimSuffix(sshKeyLabels[VMSSHKeysNames], VMSSHKeysNamesSeparator)
	return sshKeyLabels
}

func FromSSHKeySecretLabels(secretLabels map[string]string) model.SSHKeyList {
	sshKeyList := model.SSHKeyList{}
	names := strings.Split(secretLabels[VMSSHKeysNames], VMSSHKeysNamesSeparator)
	for _, name := range names {
		sshKeyList = append(sshKeyList, &model.SSHKey{
			SecretNamespace: secretLabels[VMSSHKeysSecretNamespace],
			SecretName:      name,
		})
	}
	return sshKeyList
}

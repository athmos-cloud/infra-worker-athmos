package crossplane

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"strings"
)

const (
	ManagedByLabel           = "app.kubernetes.io/managed-by"
	ManagedByValue           = "athmos"
	VMSSHKeysSecretNamespace = "vm-ssh-keys-secret-namespace"
	VMSSHKeysNamePrefix      = "vm-ssh-keys_"
	VMSSHKeysNameSeparator   = "_"
	VMPublicIPLabel          = "vm-has-public-ip"
)

func GetBaseLabels(projectID string) map[string]string {
	return map[string]string{
		ManagedByLabel:          ManagedByValue,
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
		sshKeyLabels[VMSSHKeysNamePrefix+key.Username] = key.SecretName
	}
	return sshKeyLabels
}

func FromSSHKeySecretLabels(labels map[string]string) model.SSHKeyList {
	sshKeyList := model.SSHKeyList{}
	for key, val := range labels {
		if strings.HasPrefix(key, VMSSHKeysNamePrefix) {
			splitLabel := strings.Split(val, VMSSHKeysNameSeparator)
			name := splitLabel[len(splitLabel)-1]
			sshKeyList = append(sshKeyList, &model.SSHKey{
				Username:        strings.TrimPrefix(key, VMSSHKeysNamePrefix),
				SecretNamespace: labels[VMSSHKeysSecretNamespace],
				SecretName:      name,
			})
		}
	}
	return sshKeyList
}

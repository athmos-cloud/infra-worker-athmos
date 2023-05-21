package crossplane

import v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

func GetDeletionPolicy(managed bool) v1.DeletionPolicy {
	if managed {
		return v1.DeletionDelete
	}
	return v1.DeletionOrphan
}

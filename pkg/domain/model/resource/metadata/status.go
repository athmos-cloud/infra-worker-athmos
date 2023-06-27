package metadata

import (
	cpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	corev1 "k8s.io/api/core/v1"
	"time"
)

type Status struct {
	StatusType StatusType `json:"statusType"`
	Message    string     `json:"message"`
	Date       time.Time  `json:"date"`
}

type StatusType string

const (
	StatusTypeCreateRequestSent StatusType = "CreateRequestSent"
	StatusTypeUpdateRequestSent StatusType = "UpdateRequestSent"
	StatusTypeDeleteRequestSent StatusType = "DeleteRequestSent"
	StatusTypeCreated           StatusType = "Created"
	StatusTypeUnknown           StatusType = "Unknown"
	StatusTypeError             StatusType = "Error"
	StatusTypeCreating          StatusType = "Creating"
	StatusTypeDeleting          StatusType = "Deleting"
)

func StatusFromKubernetesStatus(status []cpv1.Condition) Status {
	metaStatus := Status{}
	StatusTypeFromKubernetesStatus(&metaStatus, status)
	return metaStatus
}

func StatusTypeFromKubernetesStatus(metaStatus *Status, status []cpv1.Condition) {
	if len(status) == 0 {
		metaStatus.StatusType = StatusTypeUnknown
		return
	}
	for _, s := range status {
		if s.Status == corev1.ConditionFalse && string(s.Reason) == "ApplyFailure" {
			*metaStatus = Status{
				StatusType: StatusTypeError,
				Message:    s.Message,
				Date:       s.LastTransitionTime.Time,
			}
			return
		}
	}
	s := status[0]
	metaStatus.Message = s.Message
	metaStatus.Date = s.LastTransitionTime.Time
	if s.Reason == cpv1.ReasonAvailable {
		metaStatus.StatusType = StatusTypeCreated
		return
	}
	if s.Reason == cpv1.ReasonReconcileSuccess && s.Type == cpv1.TypeSynced {
		metaStatus.StatusType = StatusTypeCreated
		return
	}
	if s.Reason == cpv1.ReasonDeleting {
		metaStatus.StatusType = StatusTypeDeleting
		return
	}
	if s.Reason == cpv1.ReasonCreating {
		metaStatus.StatusType = StatusTypeCreating
		return
	}
	if s.Reason == cpv1.ReasonReconcileError {
		metaStatus.StatusType = StatusTypeError
		return
	}
	metaStatus.StatusType = StatusTypeUnknown
}

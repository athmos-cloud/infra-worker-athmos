package metadata

import (
	cpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
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
	s := status[0]
	return Status{
		StatusType: StatusTypeFromKubernetesStatus(status),
		Message:    s.Message,
		Date:       s.LastTransitionTime.Time,
	}
}

func StatusTypeFromKubernetesStatus(status []cpv1.Condition) StatusType {
	s := status[0]
	if s.Reason == cpv1.ReasonAvailable {
		return StatusTypeCreated
	}
	if s.Reason == cpv1.ReasonDeleting {
		return StatusTypeDeleting
	}
	if s.Reason == cpv1.ReasonCreating {
		return StatusTypeCreating
	}
	if s.Reason == cpv1.ReasonReconcileError {
		return StatusTypeError
	}
	return StatusTypeUnknown
}

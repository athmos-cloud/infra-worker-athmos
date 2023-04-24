package kubernetes

import "time"

type EventType string

const (
	Warning EventType = "warning"
)

type Reason string

const (
	PendingExternalResource  Reason = "PendingExternalResource"
	CreatingExternalResource Reason = "CreatingExternalResource"
)

type Event struct {
	Type    EventType
	Date    time.Time
	Reason  string
	Message string
}

type EventList []Event

package kubernetes

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type EventType string

const (
	Warning EventType = "warning"
)

type Reason string

const (
	PendingExternalResource  Reason = "PendingExternalResource"
	CreatingExternalResource Reason = "CreatingExternalResource"
	CannotConnectToProvider  Reason = "CannotConnectToProvider"
)

type Event struct {
	mgm.DefaultModel `bson:",inline"`
	Type             EventType `bson:"type"`
	Date             time.Time `bson:"date"`
	Reason           string    `bson:"reason"`
	Message          string    `bson:"message"`
}

func (event *Event) Equals(other Event) bool {
	return event.Type == other.Type &&
		event.Date.Equal(other.Date) &&
		event.Reason == other.Reason &&
		event.Message == other.Message
}

type EventList []Event

func (EventList *EventList) Equals(other EventList) bool {
	if len(*EventList) != len(other) {
		return false
	}
	for i, event := range *EventList {
		if !event.Equals(other[i]) {
			return false
		}
	}
	return true
}

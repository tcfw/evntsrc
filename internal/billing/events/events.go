package events

import (
	"github.com/tcfw/evntsrc/internal/utils/sysevents"
)

const (
	//BroadcastTypeCreated used when a billing customer is created
	BroadcastTypeCreated = "io.evntsrc.billing.created"
)

//Event used for generic billing svc events
type Event struct {
	*sysevents.Event
	UserID     string
	CustomerID string
}

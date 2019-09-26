package events

import "github.com/tcfw/evntsrc/internal/utils/sysevents"

const (
	//BroadcastTypeCreated used when a user is created
	BroadcastTypeCreated = "io.evntsrc.users.created"

	//BroadcastTypeActivated used when a user confirms their email
	BroadcastTypeActivated = "io.evntsrc.users.activated"

	//BroadcastTypeUpdated used when a user is updated
	BroadcastTypeUpdated = "io.evntsrc.users.updated"

	//BroadcastTypeDeleted used when a user is deleted
	BroadcastTypeDeleted = "io.evntsrc.users.deleted"

	//BroadcastTypeRecreation used when a user tries to create an account with an existing email address
	BroadcastTypeRecreation = "io.evntsrc.users.recreation"
)

//Event used for broadcasting general user svc events
type Event struct {
	*sysevents.Event
	UserID string `json:"id"`
}

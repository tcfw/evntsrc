package evntsrc

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tcfw/evntsrc/internal/websocks"
)

//Replay starts replaying events in chronological order
//Justme defaults to true if not specified
func (api *APIClient) Replay(subject string, query ReplayQuery, justme bool) error {
	cmd := &websocks.ReplayCommand{
		SubscribeCommand: &websocks.SubscribeCommand{
			InboundCommand: &websocks.InboundCommand{Ref: uuid.New().String(), Command: "replay"},
			Subject:        subject,
		},
		JustMe: justme,
		Query:  query,
	}

	api.replayPipe <- cmd

	if _, err := api.waitForResponse(cmd.Ref); err != nil && err.Error() != "OK" {
		return fmt.Errorf("Failed to start replay: %s", err.Error())
	}

	return nil
}

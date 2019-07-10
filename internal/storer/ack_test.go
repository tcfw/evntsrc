package storer

import (
	"testing"
	"time"

	"github.com/google/uuid"
	pbEvent "github.com/tcfw/evntsrc/internal/event/protos"
)

func Test_ackEvent(t *testing.T) {
	stop, err := setupTestDB(t)
	if err != nil {
		t.Error(err)
	}
	defer stop()

	proc := &eventProcessor{}
	eTime := time.Now()

	tests := []struct {
		name    string
		event   *pbEvent.Event
		wantErr bool
	}{
		{
			name: "test 1",
			event: &pbEvent.Event{
				ID:      uuid.New().String(),
				Stream:  1,
				Subject: "test",
				Source:  "test",
				Type:    "test",
				Time:    &eTime,
				Data:    []byte{},
			},
			wantErr: false,
		},
		{
			name: "test 2",
			event: &pbEvent.Event{
				ID:           uuid.New().String(),
				Stream:       1,
				Subject:      "test",
				Source:       "test",
				Type:         "test",
				Acknowledged: &eTime,
				Time:         &eTime,
				Data:         []byte{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Create event
			proc.Handle(tt.event)

			//Test func call
			if _, err := ackEvent(tt.event.Stream, tt.event.ID); (err != nil) != tt.wantErr {
				t.Errorf("ackEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

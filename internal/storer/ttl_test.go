package storer

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	pbEvent "github.com/tcfw/evntsrc/internal/event/protos"
	pb "github.com/tcfw/evntsrc/internal/storer/protos"
	pbMock "github.com/tcfw/evntsrc/internal/storer/protos/mock_protos"
)

func Test_extendTTL(t *testing.T) {
	stop, err := setupTestDB(t)
	if err != nil {
		t.Error(err)
	}
	defer stop()

	proc := &eventProcessor{}
	eTime := time.Now()

	ackTTL := time.Now().Add(20 * time.Minute)

	tests := []struct {
		name    string
		args    *pb.ExtendTTLRequest
		event   *pbEvent.Event
		wantErr bool
	}{
		{
			name: "test 1 - Add TTL to empty",
			args: &pb.ExtendTTLRequest{Stream: 1, TTLTime: &ackTTL, CurrentTTL: &ackTTL},
			event: &pbEvent.Event{
				ID:       uuid.New().String(),
				Stream:   1,
				Subject:  "test",
				Source:   "test",
				Type:     "test",
				Time:     &eTime,
				Data:     []byte{},
				Metadata: map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "test 2 - Validate existing TTL",
			args: &pb.ExtendTTLRequest{Stream: 1, TTLTime: &ackTTL, CurrentTTL: &ackTTL},
			event: &pbEvent.Event{
				ID:      uuid.New().String(),
				Stream:  1,
				Subject: "test",
				Source:  "test",
				Type:    "test",
				Time:    &eTime,
				Data:    []byte{},
				Metadata: map[string]string{
					"ttl": ackTTL.Add(-1 * time.Minute).Format(time.RFC3339),
				},
			},
			wantErr: true,
		},
		{
			name: "test 3 - No existing metadata",
			args: &pb.ExtendTTLRequest{Stream: 1, TTLTime: &ackTTL, CurrentTTL: &ackTTL},
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
			name:    "test 4 - Missing TTL time",
			wantErr: true,
			args:    &pb.ExtendTTLRequest{Stream: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.event != nil {
				//Create event
				proc.Handle(tt.event)
				tt.args.EventID = tt.event.ID
			}

			if err := extendTTL(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("extendTTL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_handleTTLQuery(t *testing.T) {
	stop, err := setupTestDB(t)
	if err != nil {
		t.Error(err)
	}
	defer stop()

	proc := &eventProcessor{}

	now := time.Now()
	pastAckTime := time.Now().Add(-1 * time.Hour)
	oldestTTL := time.Now().Add(MaxTTLDiff).Add(-1 * time.Second)

	tests := []struct {
		name        string
		req         *pb.QueryRequest
		events      []*pbEvent.Event
		wantErr     bool
		wantSend    bool
		nSentEvents int
	}{
		{
			name:     "test 1 - empty event set",
			req:      &pb.QueryRequest{Stream: 1, Query: &pb.QueryRequest_Ttl{Ttl: &pb.QueryTTLExpired{Time: &now}}},
			events:   []*pbEvent.Event{},
			wantErr:  false,
			wantSend: false,
		},
		{
			name:     "test 2 - Invalid query",
			req:      &pb.QueryRequest{Stream: 1},
			events:   []*pbEvent.Event{},
			wantErr:  true,
			wantSend: false,
		},
		{
			name:     "test 2 - Over limit",
			req:      &pb.QueryRequest{Stream: 1, Limit: 2000, Query: &pb.QueryRequest_Ttl{Ttl: &pb.QueryTTLExpired{Time: &now}}},
			events:   []*pbEvent.Event{},
			wantErr:  true,
			wantSend: false,
		},
		{
			name: "test 3 - mixed events expired and acked, but none to replay",
			req:  &pb.QueryRequest{Stream: 1, Query: &pb.QueryRequest_Ttl{Ttl: &pb.QueryTTLExpired{Time: &now}}},
			events: []*pbEvent.Event{
				{
					ID:           uuid.New().String(),
					Stream:       1,
					Subject:      "test",
					Source:       "test",
					Type:         "test",
					Time:         &now,
					Acknowledged: &pastAckTime,
					Data:         []byte{},
				},
				{
					ID:      uuid.New().String(),
					Stream:  1,
					Subject: "test",
					Source:  "test",
					Type:    "test",
					Time:    &oldestTTL,
					Data:    []byte{},
				},
			},
			wantErr:  false,
			wantSend: false,
		},
		{
			name: "test 4 - mixed events to replay for TTL",
			req:  &pb.QueryRequest{Stream: 1, Query: &pb.QueryRequest_Ttl{Ttl: &pb.QueryTTLExpired{Time: &now}}},
			events: []*pbEvent.Event{
				{
					ID:           uuid.New().String(),
					Stream:       1,
					Subject:      "test",
					Source:       "test",
					Type:         "test",
					Time:         &now,
					Acknowledged: &pastAckTime,
					Data:         []byte{},
				},
				{
					ID:      uuid.New().String(),
					Stream:  1,
					Subject: "test",
					Source:  "test",
					Type:    "test",
					Time:    &now,
					Data:    []byte{},
				},
			},
			wantErr:     false,
			wantSend:    true,
			nSentEvents: 1,
		},
		{
			name: "test 4 - all events to replay for TTL with and without prior TTL MD",
			req:  &pb.QueryRequest{Stream: 1, Query: &pb.QueryRequest_Ttl{Ttl: &pb.QueryTTLExpired{Time: &now}}},
			events: []*pbEvent.Event{
				{
					ID:      uuid.New().String(),
					Stream:  1,
					Subject: "test",
					Source:  "test",
					Type:    "test",
					Time:    &now,
					Metadata: map[string]string{
						"ttl": time.Now().Add(-5 * time.Second).Format(time.RFC3339),
					},
					Data: []byte{},
				},
				{
					ID:      uuid.New().String(),
					Stream:  1,
					Subject: "test",
					Source:  "test",
					Type:    "test",
					Time:    &now,
					Data:    []byte{},
				},
			},
			wantErr:     false,
			wantSend:    true,
			nSentEvents: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			for _, event := range tt.events {
				proc.Handle(event)
			}

			stream := pbMock.NewMockStorerService_QueryServer(ctrl)
			if tt.wantSend {
				stream.EXPECT().Send(gomock.Any()).Return(nil).Times(tt.nSentEvents)
			}

			if err := s.handleTTLQuery(tt.req, stream); (err != nil) != tt.wantErr {
				t.Errorf("server.handleTTLQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Create event
			proc.Handle(tt.event)

			tt.args.EventID = tt.event.ID

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

	type args struct {
		req    *pb.QueryRequest
		stream pb.StorerService_QueryServer
	}
	tests := []struct {
		name     string
		args     args
		events   []*pbEvent.Event
		wantErr  bool
		wantSend bool
	}{
		{
			name: "test 1 - empty event set",
			args: args{
				req:    &pb.QueryRequest{Stream: 1, Query: &pb.QueryRequest_Ttl{}},
				stream: pbMock.NewMockStorerService_QueryServer(gomock.NewController(t)),
			},
			events:   []*pbEvent.Event{},
			wantErr:  false,
			wantSend: false,
		},
		{
			name: "test 2 - Invalid query",
			args: args{
				req:    &pb.QueryRequest{Stream: 1},
				stream: pbMock.NewMockStorerService_QueryServer(gomock.NewController(t)),
			},
			events:   []*pbEvent.Event{},
			wantErr:  true,
			wantSend: false,
		},
		{
			name: "test 2 - Over limit",
			args: args{
				req:    &pb.QueryRequest{Stream: 1, Limit: 2000, Query: &pb.QueryRequest_Ttl{}},
				stream: pbMock.NewMockStorerService_QueryServer(gomock.NewController(t)),
			},
			events:   []*pbEvent.Event{},
			wantErr:  true,
			wantSend: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{}
			if tt.wantSend {
				tt.args.stream.(*pbMock.MockStorerService_QueryServer).EXPECT().Send(gomock.Any()).Return(gomock.Any())
			}
			for _, event := range tt.events {
				proc.Handle(event)
			}
			if err := s.handleTTLQuery(tt.args.req, tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("server.handleTTLQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

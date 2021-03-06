package storer

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	nats "github.com/nats-io/nats.go"

	"github.com/tcfw/evntsrc/internal/websocks"

	"github.com/stretchr/testify/assert"

	pbEvent "github.com/tcfw/evntsrc/internal/event/protos"

	"github.com/cockroachdb/cockroach-go/testserver"
	natsTest "github.com/nats-io/nats-server/v2/test"
)

func setupTestDB(t *testing.T) (func(), error) {
	db, stop := testserver.NewDBForTest(t)

	pgdb = db

	if err := createUpdateTables(pgdb); err != nil {
		return nil, err
	}

	return stop, nil
}

func TestStartMonitor(t *testing.T) {
	s := natsTest.RunDefaultServer()
	s.Start()
	defer s.Shutdown()

	go StartMonitor(fmt.Sprintf("nats://%s", s.Addr().String()))

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, uint32(3), s.NumSubscriptions())
}

func TestMonitorUserStreams(t *testing.T) {
	s := natsTest.RunDefaultServer()
	s.Start()
	defer s.Shutdown()

	natsConn, _ = nats.Connect(fmt.Sprintf("nats://%s", s.Addr().String()))

	monitorUserStreams()

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, uint32(1), s.NumSubscriptions())
}

func TestEventHandle(t *testing.T) {
	stop, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer stop()

	proc := &eventProcessor{}
	eTime := time.Now()
	event := &pbEvent.Event{
		ID:      uuid.New().String(),
		Stream:  1,
		Subject: "test",
		Source:  "test",
		Type:    "test",
		Time:    &eTime,
		Data:    []byte{},
	}

	assert.NotPanics(t, func() {
		proc.Handle(event)
	})
}

func TestDoReply(t *testing.T) {

	s := natsTest.RunDefaultServer()
	s.Start()
	defer s.Shutdown()

	var err error
	natsConn, _ = nats.Connect(fmt.Sprintf("nats://%s", s.Addr().String()))
	if !natsConn.IsConnected() {
		t.Fatal("nats not connected")
	}

	stop, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer stop()

	reply := make(chan []byte, 10)
	errCh := make(chan error, 1)

	eTime := time.Now()
	event := &pbEvent.Event{
		ID:      uuid.New().String(),
		Stream:  1,
		Subject: "test",
		Source:  "test",
		Type:    "test",
		Time:    &eTime,
		Data:    []byte{},
	}

	proc := &eventProcessor{}
	proc.Handle(event)

	command := &websocks.ReplayCommand{Stream: 1, SubscribeCommand: &websocks.SubscribeCommand{Subject: "test", InboundCommand: &websocks.InboundCommand{Ref: "test"}}}

	doReplay(command, reply, errCh)

	replyMsg := <-reply
	err = <-errCh

	assert.NoError(t, err)
	assert.Equal(t, []byte("OK"), replyMsg)
}

func TestBuildBaseQuery(t *testing.T) {
	command := &websocks.ReplayCommand{
		Stream: 1,
	}

	query, params := buildBaseQuery(command)

	assert.Equal(t, `FROM event_store.events WHERE stream = $1`, query)
	assert.Equal(t, 1, len(params))

	sTime := time.Now()

	command = &websocks.ReplayCommand{
		Stream: 1,
		Query: websocks.ReplayRange{
			StartTime: &sTime,
		},
	}

	query, params = buildBaseQuery(command)

	assert.Equal(t, `FROM event_store.events WHERE stream = $1 AND time >= $2`, query)
	assert.Equal(t, 2, len(params))

	eTime := time.Now()

	command = &websocks.ReplayCommand{
		Stream: 1,
		Query: websocks.ReplayRange{
			StartTime: &sTime,
			EndTime:   &eTime,
		},
	}

	query, params = buildBaseQuery(command)

	assert.Equal(t, `FROM event_store.events WHERE stream = $1 AND time >= $2 AND time <= $3`, query)
	assert.Equal(t, 3, len(params))
}

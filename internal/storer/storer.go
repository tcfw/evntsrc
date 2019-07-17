package storer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	pbEvent "github.com/tcfw/evntsrc/internal/event/protos"
	pbStorer "github.com/tcfw/evntsrc/internal/storer/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
	"github.com/tcfw/evntsrc/internal/websocks"
	"github.com/tcfw/go-queue"
	"google.golang.org/grpc"
)

var (
	pgdb *sql.DB

	replayDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "storer_replay_request_duration_seconds",
		Help:    "Histogram for the runtime of a replay function.",
		Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
	})

	replayEventCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "storer_replay_event_count",
		Help: "Counter for replay events to NATS",
	}, []string{"stream"})

	storeCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "storer_store_request_count",
		Help: "Counter for store requests from NATS events",
	}, []string{"stream"})
)

//Start inits required processes
func Start(nats string, port int) {
	dbURL, ok := os.LookupEnv("PGDB_HOST")
	if !ok {
		log.Fatal("Failed to connect to DB")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pgdb = db

	if err := createUpdateTables(pgdb); err != nil {
		log.Fatal(err)
	}

	if viper.GetBool("migrate") {
		os.Exit(0)
		return
	}

	go RegisterMetrics()
	go StartGRPC(port)

	StartMonitor(nats)
}

//RegisterMetrics registers metrics with prometheus
func RegisterMetrics() {
	prometheus.MustRegister(replayDurations)
	prometheus.MustRegister(storeCount)
	prometheus.MustRegister(replayEventCount)

	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":80", nil))
}

//StartGRPC starts web http and grpc endpoint
func StartGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(tracing.GRPCServerOptions()...)

	srv, err := newServer()
	if err != nil {
		log.Fatal(err)
	}

	pbStorer.RegisterStorerServiceServer(grpcServer, srv)

	log.Println("Starting gRPC server...")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("GRPC: %s", err.Error())
	}
}

//StartMonitor subscripts to all user channels
func StartMonitor(nats string) {
	connectNats(nats)
	defer natsConn.Close()

	monitorUserStreams()
	monitorReplayRequests()

	//Wait forever
	select {}
}

type eventProcessor struct{}

func (ep *eventProcessor) Handle(job interface{}) {
	usrEvent := job.(*pbEvent.Event)

	if isReplay, ok := usrEvent.Metadata["replay"]; ok && isReplay == "true" {
		return
	}

	if _, ok := usrEvent.Metadata["forwarded"]; ok {
		return
	}

	if isNonPersistent, ok := usrEvent.Metadata["non-persistent"]; ok && isNonPersistent == "true" {
		return
	}

	metadataJSON, err := json.Marshal(usrEvent.Metadata)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := pgdb.Begin()
	if err != nil {
		panic(err.Error)
	}

	md := string(metadataJSON)
	if len(metadataJSON) == 0 || len(usrEvent.Metadata) == 0 {
		md = "{}"
	}

	if _, err := tx.Exec(
		`INSERT INTO event_store.events VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		usrEvent.ID,
		usrEvent.Stream,
		usrEvent.Time,
		usrEvent.Type,
		usrEvent.TypeVersion,
		usrEvent.CEVersion,
		usrEvent.Source,
		usrEvent.Subject,
		usrEvent.Acknowledged,
		md,
		usrEvent.ContentType,
		usrEvent.Data,
	); err != nil {
		log.Fatal(err)
		panic(err.Error)
	}

	tx.Commit()

	storeCount.With(prometheus.Labels{"stream": fmt.Sprintf("%d", usrEvent.Stream)}).Inc()

}

func monitorUserStreams() {
	log.Println("Watching for user streams...")

	dispatcher := queue.NewDispatcher(&eventProcessor{})
	dispatcher.MaxWorkers = 20
	dispatcher.Run()

	inbc := make(chan *pbEvent.Event, 1000)

	go func() {
		for {
			event := <-inbc
			dispatcher.Queue(event)
		}
	}()

	natsConn.QueueSubscribe("_USER.>", "storers", func(m *nats.Msg) {
		event := &pbEvent.Event{}
		err := proto.Unmarshal(m.Data, event)
		if err != nil {
			log.Printf("%s\n", err.Error())
			return
		}

		inbc <- event
	})
}

func monitorReplayRequests() {
	log.Println("Watching for replay requests...")
	natsConn.QueueSubscribe("replay.broadcast", "replayers", func(m *nats.Msg) {
		timer := prometheus.NewTimer(replayDurations)
		defer timer.ObserveDuration()

		command := &websocks.ReplayCommand{}
		json.Unmarshal(m.Data, command)

		reply := make(chan []byte, 10)
		errCh := make(chan error)

		go doReplay(command, reply, errCh)

		select {
		case msg := <-reply:
			natsConn.Publish(m.Reply, msg)
		case err := <-errCh:
			natsConn.Publish(m.Reply, []byte(fmt.Sprintf("failed: %s", err.Error())))
			log.Print(err.Error())
		}
	})
}

func doReplay(command *websocks.ReplayCommand, reply chan []byte, errCh chan error) {
	defer func(reply chan []byte, errCh chan error) {
		close(reply)
		close(errCh)
	}(reply, errCh)

	if command.JustMe && command.Dest == "" {
		errCh <- fmt.Errorf("no socket dest set")
		return
	}

	qSelector, params := buildBaseQuery(command)

	count, err := countEvents(qSelector, params)
	if err != nil {
		errCh <- err
		return
	}
	if count == 0 {
		reply <- []byte("no events")
		return
	}

	rD, err := pgdb.Query(`SELECT * `+qSelector, params...)
	if err != nil {
		errCh <- fmt.Errorf("sqld: %s", err.Error())
		return
	}
	defer func() {
		rD.Close()
	}()

	reply <- []byte("OK")

	for rD.Next() {
		event, err := scanEvent(rD)
		if err != nil {
			errCh <- err
		}
		if command.Query.EndID != "" && event.ID == command.Query.EndID {
			break
		}

		bytes, err := proto.Marshal(event)
		if err != nil {
			errCh <- fmt.Errorf("failed replay proto marshal: %s", err.Error())
			return
		}
		if command.JustMe {
			dest := fmt.Sprintf("_CONN.%s", command.Dest)
			err = natsConn.Publish(dest, bytes)
			if err != nil {
				errCh <- fmt.Errorf("natspub: %s", err.Error())
				return
			}
		} else {
			dest := fmt.Sprintf("_USER.%d.%s", command.Stream, command.Subject)
			err = natsConn.Publish(dest, bytes)
			if err != nil {
				errCh <- fmt.Errorf("natspub: %s", err.Error())
				return
			}
		}
		replayEventCount.With(prometheus.Labels{"stream": fmt.Sprintf("%d", command.Stream)}).Inc()
	}
}

func scanEvent(rD *sql.Rows) (*pbEvent.Event, error) {
	event := &pbEvent.Event{
		Metadata: map[string]string{},
	}
	var mdString []byte

	err := rD.Scan(&event.ID,
		&event.Stream,
		&event.Time,
		&event.Type,
		&event.TypeVersion,
		&event.CEVersion,
		&event.Source,
		&event.Subject,
		&event.Acknowledged,
		&mdString,
		&event.ContentType,
		&event.Data,
	)
	if err != nil {
		return nil, fmt.Errorf("sqld: %s", err.Error())
	}

	err = json.Unmarshal(mdString, &event.Metadata)
	if err != nil {
		return nil, fmt.Errorf("sqld md: %s", err.Error())
	}
	if event.Metadata == nil {
		event.Metadata = map[string]string{}
	}

	event.Metadata["replay"] = "true"

	return event, nil
}

func buildBaseQuery(command *websocks.ReplayCommand) (string, []interface{}) {
	qSelector := `FROM event_store.events WHERE stream = $1`

	params := []interface{}{
		command.Stream,
	}

	if command.Query.StartTime != nil && command.Query.EndTime == nil {
		qSelector += ` AND time >= $2`
		params = append(params, command.Query.StartTime)
	} else if command.Query.StartTime != nil && command.Query.EndTime != nil {
		qSelector += ` AND time >= $2 AND time <= $3`
		params = append(params, command.Query.StartTime)
		params = append(params, command.Query.EndTime)
	}

	// qSelector += ` GROUP BY id, time ORDER BY time;`
	return qSelector, params
}

func countEvents(qSelector string, params []interface{}) (int, error) {
	count := 0

	rC, err := pgdb.Query(`SELECT COUNT(id) `+qSelector, params...)
	if err != nil {
		return count, fmt.Errorf("sqlc: %s", err.Error())
	}

	rC.Next()

	if err := rC.Scan(&count); err != nil {
		return count, fmt.Errorf("sqlc: %s", err.Error())
	}
	rC.Close()

	return count, nil
}

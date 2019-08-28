package storer

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	pbEvent "github.com/tcfw/evntsrc/internal/event/protos"
	pbStorer "github.com/tcfw/evntsrc/internal/storer/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
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
	monitorInternalStreams()

	//Wait forever
	select {}
}

type eventProcessor struct{}

func (ep *eventProcessor) Handle(job interface{}) {
	usrEvent := job.(*pbEvent.Event)

	err := storeEvent(usrEvent, pgdb)
	if err != nil {
		log.Println(err)
	}

	go ep.ackPub(usrEvent)
}

func (ep *eventProcessor) ackPub(event *pbEvent.Event) {
	src, ok := event.Metadata["_cid"]
	if !ok {
		return
	}

	now := time.Now()
	ack := &pbEvent.Event{Time: &now, ID: event.ID, Stream: 0, Subject: "puback"}
	bytes, _ := proto.Marshal(ack)

	natsConn.Publish(fmt.Sprintf("_CONN.%s", src), bytes)
}

//monitorUserStreams watches for all user channels typically prefixed with _USER
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

//monitorInternalStreams watches for all internal events prefixed with _INTERNAL
func monitorInternalStreams() {
	log.Println("Watching for internal streams...")

	dispatcher := queue.NewDispatcher(&eventProcessor{})
	dispatcher.MaxWorkers = 10
	dispatcher.Run()

	inbc := make(chan *pbEvent.Event, 100)

	go func() {
		for {
			event := <-inbc
			dispatcher.Queue(event)
		}
	}()

	natsConn.QueueSubscribe("_INTERNAL.>", "storers", func(m *nats.Msg) {
		event := &pbEvent.Event{}
		err := proto.Unmarshal(m.Data, event)
		if err != nil {
			log.Printf("%s\n", err.Error())
			return
		}

		inbc <- event
	})
}

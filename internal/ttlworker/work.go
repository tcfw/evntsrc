package ttlworker

import (
	"log"
	"os"
	"sync"

	storerPB "github.com/tcfw/evntsrc/internal/storer/protos"
	ttlschedulerPB "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
	"google.golang.org/grpc"
)

//Worker primary ttl worker
type Worker struct {
	schedulerConn *grpc.ClientConn
	schedulerCli  ttlschedulerPB.TTLSchedulerClient
	storerConn    *grpc.ClientConn
	storerCli     storerPB.StorerServiceClient
	bindings      []*ttlschedulerPB.Binding
	bindingMu     sync.RWMutex
	stopWatching  chan struct{}
}

//NewWorker constructs a worker with a connection to the ttlscheduler for fetching
//allocated streams
func NewWorker(port int) (*Worker, error) {
	var schedulerEndpoint string
	var storerEndpoint string

	//Scheduler conn
	if schedulerEndpoint = os.Getenv("SCHEDULER_ENDPOINT"); schedulerEndpoint == "" {
		schedulerEndpoint = "ttlscheduler:443"
	}

	schedulerConn, err := grpc.Dial(schedulerEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	schedulerClient := ttlschedulerPB.NewTTLSchedulerClient(schedulerConn)

	//Storer conn
	if storerEndpoint = os.Getenv("STORER_ENDPOINT"); storerEndpoint == "" {
		storerEndpoint = "storer:443"
	}

	storerConn, err := grpc.Dial(storerEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	storerClient := storerPB.NewStorerServiceClient(storerConn)

	return &Worker{
		grpcPort:      port,
		schedulerConn: schedulerConn,
		schedulerCli:  schedulerClient,
		storerConn:    storerConn,
		storerCli:     storerClient,
		bindings:      []*ttlschedulerPB.Binding{},
	}, nil
}

//StartAndWait starts the worker and waits forever
func (w *Worker) StartAndWait() {
	w.Start()
	select {}
}

//Start starts monitoring ttlscheduler allocations and enacts ttl replays based on bindings
func (w *Worker) Start() {
	log.Println("Starting...")
	go w.watchBindings()
}

//Stop closes the GRPC connection to the scheduler
//and stops timed actions like ttl replays
func (w *Worker) Stop() {
	w.stopWatching <- struct{}{}
	w.schedulerConn.Close()
}

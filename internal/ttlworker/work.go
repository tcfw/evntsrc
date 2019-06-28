package ttlworker

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	ttlschedulerPB "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
	"google.golang.org/grpc"
)

//Worker primary ttl worker
type Worker struct {
	grpcPort      int //TODO(tcfw): Currently unused
	schedulerConn *grpc.ClientConn
	schedulerCli  ttlschedulerPB.TTLSchedulerClient
	bindings      []*ttlschedulerPB.Binding
	bindingMu     sync.RWMutex
	stopWatching  chan struct{}
}

//NewWorker constructs a worker with a connection to the ttlscheduler for fetching
//allocated streams
func NewWorker(port int) (*Worker, error) {
	var schedulerEndpoint string
	if schedulerEndpoint = os.Getenv("SCHEDULER_ENDPOINT"); schedulerEndpoint == "" {
		schedulerEndpoint = "ttlscheduler.default:443"
	}

	conn, err := grpc.Dial(schedulerEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	schedulerClient := ttlschedulerPB.NewTTLSchedulerClient(conn)

	return &Worker{
		grpcPort:      port,
		schedulerConn: conn,
		schedulerCli:  schedulerClient,
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

//watchBindings calls fetch streams every minute or stops working on stopWatching
func (w *Worker) watchBindings() {
	watchTicker := time.NewTicker(30 * time.Second)

	log.Println("Monitoring bindings")

	//Once off to seed
	go w.fetchBindings()

	for {
		select {
		case <-watchTicker.C:
			w.fetchBindings()
		case <-w.stopWatching:
			watchTicker.Stop()
			return
		}
	}
}

//fetchBindings calls the ttlscheduler to fetch bindings which have been allocated to this worker
//based on the hostname as the node ID
func (w *Worker) fetchBindings() error {
	id := w.identifier()

	bindingsResp, err := w.schedulerCli.NodeBindings(context.Background(), &ttlschedulerPB.NodeBindingRequest{Node: &ttlschedulerPB.Node{Id: id}})
	if err != nil {
		return err
	}

	w.bindingMu.Lock()
	w.bindings = bindingsResp.Bindings
	w.bindingMu.Unlock()

	return nil
}

func (w *Worker) identifier() string {
	hostname, _ := os.Hostname()
	return hostname
}

//Stop closes the GRPC connection to the scheduler
//and stops timed actions like ttl replays
func (w *Worker) Stop() {
	w.stopWatching <- struct{}{}
	w.schedulerConn.Close()
}

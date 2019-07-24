package ttlworker

import (
	"context"
	"log"
	"os"
	"time"

	ttlschedulerPB "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
)

//watchBindings calls fetch streams every minute or stops working on stopWatching
func (w *Worker) watchBindings() {
	watchTicker := time.NewTicker(30 * time.Second)
	streamTicker := time.NewTicker(30 * time.Second)

	log.Println("Monitoring bindings")

	//Once off to seed
	go w.fetchBindings()

	for {
		select {
		case <-watchTicker.C:
			go func() {
				if err := w.fetchBindings(); err != nil {
					log.Println(err)
				}
			}()
		case <-streamTicker.C:
			go func() {
				if err := w.processBindings(); err != nil {
					log.Println(err)
				}
			}()
		case <-w.stopWatching:
			watchTicker.Stop()
			streamTicker.Stop()
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

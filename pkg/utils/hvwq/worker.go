package hvwq

import (
	"reflect"
	"time"
)

//WorkProcessor should be able to process jobs that come in
type WorkProcessor interface {
	Handle(job interface{})
}

//Worker proccesses jobs via the processor interface
type Worker struct {
	WorkerPool chan chan interface{}
	JobChannel chan interface{}
	quit       chan bool
	Processor  WorkProcessor
}

//NewWorker bootstraps the job and quit channels
func NewWorker(workerPool chan chan interface{}) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan interface{}),
		quit:       make(chan bool, 1),
	}
}

//Ready attaches the worker to the pool
func (w Worker) Ready() {
	if w.WorkerPool != nil {
		w.WorkerPool <- w.JobChannel
	}
}

//Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.JobChannel = nil
		w.quit <- true
	}()
}

//Start beings processing jobs from the worker's job queue
func (w Worker) Start() {
	go func() {
		for w.JobChannel != nil {
			// register the current worker into the worker queue.
			w.Ready()

			select {
			case job, ok := <-w.JobChannel:
				// we have received a work request.
				if ok {
					start := time.Now()

					w.Processor.Handle(job)

					duration := time.Since(start)
					jobsDurationHistogram.WithLabelValues(reflect.TypeOf(job).String()).Observe(duration.Seconds())
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

package hvwq

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

var jobsDurationHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "jobs_duration_seconds",
		Help:    "Jobs duration distribution",
		Buckets: []float64{1, 2, 5, 10, 20, 60},
	},
	[]string{"job_type"},
)

//Dispatcher holds the worker pool and delegates jobs to those workers
type Dispatcher struct {
	jobQueue   chan interface{}
	WorkerPool chan chan interface{}
	MaxWorkers int
	Processor  WorkProcessor
	Workers    []Worker
}

//NewDispatcher creates a new dispatcher with respective processor
func NewDispatcher(processor WorkProcessor) *Dispatcher {
	return &Dispatcher{
		jobQueue:   make(chan interface{}),
		WorkerPool: make(chan chan interface{}),
		Processor:  processor,
		Workers:    []Worker{},
	}
}

//Run registers the workers and allocation of jobs
func (d *Dispatcher) Run() {

	if d.MaxWorkers == 0 {
		d.MaxWorkers = runtime.NumCPU()
	}

	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Processor = d.Processor
		d.Workers = append(d.Workers, worker)
		worker.Start()
	}

	go d.dispatch()
}

//Queue sends a job into to be processed
func (d *Dispatcher) Queue(job interface{}) {
	d.jobQueue <- job
}

func (d *Dispatcher) dispatch() {
	for d.jobQueue != nil && d.WorkerPool != nil {
		if job, ok := <-d.jobQueue; ok {
			go func(job interface{}) {
				if JobChannel, ok := <-d.WorkerPool; ok {
					JobChannel <- job
				}
			}(job)
		}
	}
}

//Stop signals all queue workers to quit
func (d *Dispatcher) Stop() {
	d.jobQueue = nil
	d.WorkerPool = nil
	for _, worker := range d.Workers {
		worker.Stop()
	}
}

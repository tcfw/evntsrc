package hvwq

import (
	"fmt"
	"sync"
	"testing"
)

type Job struct {
	msg string
}

func (j *Job) getMessage() string {
	return j.msg
}

type Processor struct{}

func (p Processor) Handle(job interface{}) {
	job.(struct {
		wait *sync.WaitGroup
	}).wait.Done()
}

func DispatcherNWorkers(b *testing.B, workers int) {
	dispatcher := NewDispatcher(&Processor{})
	dispatcher.MaxWorkers = workers
	dispatcher.Run()

	var wg sync.WaitGroup
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		dispatcher.Queue(struct {
			wait *sync.WaitGroup
		}{
			wait: &wg,
		})
	}

	wg.Wait()

	dispatcher.Stop()
}

func BenchmarkDispatcher(b *testing.B) {
	for n := 1; n < 10; n += 2 {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			DispatcherNWorkers(b, n)
		})
	}
}

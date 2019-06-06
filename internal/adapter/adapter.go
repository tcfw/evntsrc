package adapter

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/tcfw/evntsrc/internal/adapter/protos"
	v8 "gopkg.in/augustoroman/v8.v1"
)

//Server core struct
type Server struct {
	mu         sync.Mutex
	v8Pool     chan *v8.Context
	v8PoolStop chan struct{}
}

//NewServer creates a new struct to interface the streams server
func NewServer() *Server {
	return &Server{}
}

//StartPools starts all engine pools
func (s *Server) StartPools() {
	s.StartV8Pool()
	fmt.Println("Waiting for V8 Pool to populate at least 5 isolated contexts")
	for {
		if len(s.v8Pool) < 5 {
			fmt.Printf(".")
			time.Sleep(1 * time.Second)
		} else {
			fmt.Printf("\n")
			return
		}
	}
}

//StopPools stops and clears all engine pools
func (s *Server) StopPools() {
	s.StopV8Pool()
}

//StartV8Pool creates a continuous stream of ready V8 isolated contexts
func (s *Server) StartV8Pool() {
	s.v8Pool = make(chan *v8.Context, 20)
	s.v8PoolStop = make(chan struct{})
	go func() {
		for {
			select {
			case <-s.v8PoolStop:
				return
			default:
				s.v8Pool <- v8.NewIsolate().NewContext()
			}
		}
	}()
}

//StopV8Pool clears all V8 isolated contexts from the pool
func (s *Server) StopV8Pool() {
	close(s.v8PoolStop)
	close(s.v8Pool)
}

//Execute mutates events based on adapters
func (s *Server) Execute(ctx context.Context, request *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	switch engineType := request.Adapter.Engine; engineType {
	case pb.Adapter_V8:
		event, log, err := RunV8Adapter(s, request.Adapter, request.Event)
		if err != nil {
			return nil, err
		}

		return &pb.ExecuteResponse{Event: event, Log: log}, err
	default:
		return nil, fmt.Errorf("Unknown engine type: %s", engineType)
	}
}

package ttlscheduler

import (
	"context"
	"errors"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/spf13/viper"

	pb "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
)

type node struct {
	ID int `json:"id,omitempty"`
}

type stream struct {
	ID      int32 `json:"id,omitempty"`
	MsgRate int   `json:"msg_rate,omitempty"`
}

type binding struct {
	Stream *stream `json:"stream_id,omitempty"`
	Node   *node   `json:"node,omitempty"`
}

//NodeFetcher builds up the node list form a source e.g. kubernetes pods
type NodeFetcher interface {
	GetNodes() ([]*node, error)
}

//StreamFetcher provides intereface to fetch stream lists
type StreamFetcher interface {
	GetStreams() ([]*stream, error)
}

//Scheduler basic structure of a scheduler
type Scheduler interface {
	NodeBindings(context.Context, *pb.NodeBindingRequest) (*pb.NodeBindingResponse, error)
	BindStream(*stream) (*binding, error)
	Optimise() error
	Observe() error
}

type basicScheduler struct {
	nodes    map[int]*node
	streams  []*stream
	bindings []*binding
	nf       NodeFetcher
	sf       StreamFetcher
	once     bool

	observeInternal time.Duration
	lock            sync.RWMutex
}

func (s *basicScheduler) NodeBindings(ctx context.Context, req *pb.NodeBindingRequest) (*pb.NodeBindingResponse, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	nodeBindings := []*pb.Binding{}

	for _, bind := range s.bindings {
		if bind.Node == nil {
			continue
		}
		if bind.Node.ID == int(req.Node.Id) {
			binding := &pb.Binding{
				Stream: &pb.Stream{Id: bind.Stream.ID, MsgRate: int64(bind.Stream.MsgRate)},
				Node:   &pb.Node{Id: int32(bind.Node.ID)},
			}
			nodeBindings = append(nodeBindings, binding)
		}
	}

	return &pb.NodeBindingResponse{Bindings: nodeBindings}, nil
}

func (s *basicScheduler) GetNodes() map[int]*node {
	return s.nodes
}

//BindStream allocates a stream to a node worker in an "optimal" pattern
func (s *basicScheduler) BindStream(stream *stream) (*binding, error) {
	nCount := len(s.nodes)

	if nCount == 0 {
		return nil, errors.New("No nodes to schedule on")
	}
	if nCount == 1 {
		binding := &binding{Stream: stream, Node: s.nodes[int(reflect.ValueOf(s.nodes).MapKeys()[0].Int())]}
		return binding, nil
	}

	nScores := map[int]int{}
	for i, node := range s.nodes {
		nScores[i] = s.nodeScore(*node)
	}

	lowestScore, lowestNode := -1, -1
	for node, score := range nScores {
		if lowestNode == -1 || score < lowestScore {
			lowestNode = node
			lowestScore = score
		}
	}

	binding := &binding{Stream: stream, Node: s.nodes[lowestNode]}
	return binding, nil
}

//NodeScore calculates a load score for a particular node
func (s *basicScheduler) nodeScore(node node) int {
	score := 0
	nodeBindings := []*binding{}

	for _, bind := range s.bindings {
		if bind.Node == nil {
			continue
		}
		if bind.Node.ID == node.ID {
			nodeBindings = append(nodeBindings, bind)
		}
	}

	for _, binding := range nodeBindings {
		score += 10 + binding.Stream.MsgRate
	}

	return score
}

//Optimise @TODO
func (s *basicScheduler) Optimise() error {
	return nil
}

func (s *basicScheduler) Observe() error {
	i := 0
	for {
		nodes, err := s.nf.GetNodes()
		if err != nil {
			return err
		}

		streams, err := s.sf.GetStreams()
		if err != nil {
			return err
		}

		s.lock.Lock()

		s.observeNodes(nodes)
		s.observeStreams(streams)

		if i%6 == 0 && viper.GetBool("verbose") == true {
			log.Printf("Streams: %d; Nodes: %d; Bindings: %d\n", len(s.streams), len(s.nodes), len(s.bindings))
		}
		i++

		s.lock.Unlock()

		if s.once {
			return nil
		}

		time.Sleep(s.observeInternal)
	}
}

func (s *basicScheduler) observeNodes(nNodes []*node) {
	if len(nNodes) == 0 {
		return
	}

	addedNodes, deletedNodes := s.nodeDiff(nNodes)
	for id, added := range addedNodes {
		s.nodes[id] = added
	}

	for id := range deletedNodes {
		reschedule := []*stream{}
		for _, binding := range s.bindings {
			if binding.Node.ID == id {
				reschedule = append(reschedule, binding.Stream)
				binding.Node = nil
			}
		}

		delete(s.nodes, id)

		for _, resStream := range reschedule {
			binding, _ := s.BindStream(&stream{resStream.ID, 0})
			s.bindings = append(s.bindings, binding)
		}
	}

	final := []*binding{}
	for _, binding := range s.bindings {
		if binding.Node != nil {
			final = append(final, binding)
		}
	}

	s.bindings = final
}

func (s *basicScheduler) observeStreams(nStreams []*stream) {
	if len(nStreams) == 0 {
		return
	}
	addedStreams, deletedStreams := s.streamDiff(nStreams)
	for id, aStream := range addedStreams {
		binding, _ := s.BindStream(&stream{id, 0})
		s.bindings = append(s.bindings, binding)
		s.streams = append(s.streams, aStream)
	}

	if len(deletedStreams) > 0 {
		final := []*binding{}
		for id := range deletedStreams {
			for _, binding := range s.bindings {
				if binding.Stream.ID != id {
					final = append(final, binding)
				}
			}

			index := -1
			for i, stream := range s.streams {
				if stream.ID == id {
					index = i
					break
				}
			}
			if index != -1 {
				s.streams = append(s.streams[:index], s.streams[index+1:]...)
			}
		}

		s.bindings = final
	}
}

//TODO merge streamDiff and nodeDiff
func (s *basicScheduler) streamDiff(nStreams []*stream) (map[int32]*stream, map[int32]*stream) {
	added := map[int32]*stream{}
	deleted := map[int32]*stream{}
	same := map[int32]*stream{}

	streamMap := map[int32]*stream{}
	for _, sStream := range s.streams {
		streamMap[sStream.ID] = sStream
	}

	//Find new and same
	for _, nStream := range nStreams {
		if _, ok := streamMap[nStream.ID]; ok {
			same[nStream.ID] = nStream
		} else {
			added[nStream.ID] = nStream
		}
	}

	//Find deleted
	for oSID, oStream := range streamMap {
		_, inA := added[oSID]
		_, inS := same[oSID]
		if !inA && !inS {
			deleted[oSID] = oStream
		}
	}

	return added, deleted
}

func (s *basicScheduler) nodeDiff(nNodes []*node) (map[int]*node, map[int]*node) {
	added := map[int]*node{}
	deleted := map[int]*node{}
	same := map[int]*node{}

	//Find new and same
	for _, nNode := range nNodes {
		if _, ok := s.nodes[nNode.ID]; ok {
			same[nNode.ID] = nNode
		} else {
			added[nNode.ID] = nNode
		}
	}

	//Find deleted
	for oNID, oNode := range s.nodes {
		_, inA := added[oNID]
		_, inS := same[oNID]
		if !inA && !inS {
			deleted[oNID] = oNode
		}
	}

	return added, deleted
}

//NewScheduler constructs a new scheduler which can assign streams work nodes
func NewScheduler(nf NodeFetcher, sf StreamFetcher) Scheduler {
	return &basicScheduler{
		nodes:           map[int]*node{},
		streams:         []*stream{},
		bindings:        []*binding{},
		nf:              nf,
		sf:              sf,
		observeInternal: 5 * time.Second,
	}
}

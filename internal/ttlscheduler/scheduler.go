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

//NodeFetcher builds up the node list form a source e.g. kubernetes pods
type NodeFetcher interface {
	GetNodes() ([]*pb.Node, error)
}

//StreamFetcher provides intereface to fetch stream lists
type StreamFetcher interface {
	GetStreams() ([]*pb.Stream, error)
}

//Scheduler basic structure of a scheduler
type Scheduler interface {
	NodeBindings(context.Context, *pb.NodeBindingRequest) (*pb.NodeBindingResponse, error)
	BindStream(*pb.Stream) (*pb.Binding, error)
	Optimise() error
	Observe() error
}

type basicScheduler struct {
	nodes    map[string]*pb.Node
	streams  []*pb.Stream
	bindings []*pb.Binding
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
		if bind.Node.Id == req.Node.Id {
			binding := &pb.Binding{
				Stream: &pb.Stream{Id: bind.Stream.Id, MsgRate: int64(bind.Stream.MsgRate)},
				Node:   &pb.Node{Id: bind.Node.Id},
			}
			nodeBindings = append(nodeBindings, binding)
		}
	}

	return &pb.NodeBindingResponse{Bindings: nodeBindings}, nil
}

func (s *basicScheduler) GetNodes() map[string]*pb.Node {
	return s.nodes
}

//BindStream allocates a stream to a node worker in an "optimal" pattern
func (s *basicScheduler) BindStream(stream *pb.Stream) (*pb.Binding, error) {
	nCount := len(s.nodes)

	if nCount == 0 {
		return nil, errors.New("No nodes to schedule on")
	}
	if nCount == 1 {
		binding := &pb.Binding{Stream: stream, Node: s.nodes[reflect.ValueOf(s.nodes).MapKeys()[0].String()]}
		return binding, nil
	}

	nScores := map[string]int64{}
	for i, node := range s.nodes {
		nScores[i] = s.nodeScore(node)
	}

	lowestScore, lowestNode := int64(-1), ""
	for node, score := range nScores {
		if lowestNode == "" || score < lowestScore {
			lowestNode = node
			lowestScore = score
		}
	}

	binding := &pb.Binding{Stream: stream, Node: s.nodes[lowestNode]}
	return binding, nil
}

//NodeScore calculates a load score for a particular node
func (s *basicScheduler) nodeScore(node *pb.Node) int64 {
	score := int64(0)
	nodeBindings := []*pb.Binding{}

	for _, bind := range s.bindings {
		if bind.Node == nil {
			continue
		}
		if bind.Node.Id == node.Id {
			nodeBindings = append(nodeBindings, bind)
		}
	}

	for _, binding := range nodeBindings {
		score += 10 + binding.Stream.MsgRate
	}

	return score
}

//Optimise reviews all allocations and attempts to optimise
//to even the load between each node available
func (s *basicScheduler) Optimise() error {
	//TODO(tcfw)
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

func (s *basicScheduler) observeNodes(nNodes []*pb.Node) {
	if len(nNodes) == 0 {
		return
	}

	addedNodes, deletedNodes := s.nodeDiff(nNodes)
	for id, added := range addedNodes {
		s.nodes[id] = added
	}

	for id := range deletedNodes {
		reschedule := []*pb.Stream{}
		for _, binding := range s.bindings {
			if binding.Node.Id == id {
				reschedule = append(reschedule, binding.Stream)
				binding.Node = nil
			}
		}

		delete(s.nodes, id)

		for _, resStream := range reschedule {
			binding, _ := s.BindStream(&pb.Stream{Id: resStream.Id, MsgRate: 0})
			s.bindings = append(s.bindings, binding)
		}
	}

	final := []*pb.Binding{}
	for _, binding := range s.bindings {
		if binding.Node != nil {
			final = append(final, binding)
		}
	}

	s.bindings = final
}

func (s *basicScheduler) observeStreams(nStreams []*pb.Stream) {
	if len(nStreams) == 0 {
		return
	}
	addedStreams, deletedStreams := s.streamDiff(nStreams)
	for id, aStream := range addedStreams {
		binding, err := s.BindStream(&pb.Stream{Id: id, MsgRate: 0})
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		s.bindings = append(s.bindings, binding)
		s.streams = append(s.streams, aStream)
	}

	if len(deletedStreams) > 0 {
		final := []*pb.Binding{}
		for id := range deletedStreams {
			for _, binding := range s.bindings {
				if binding.Stream.Id != id {
					final = append(final, binding)
				}
			}

			index := -1
			for i, stream := range s.streams {
				if stream.Id == id {
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

func (s *basicScheduler) streamDiff(nStreams []*pb.Stream) (map[int32]*pb.Stream, map[int32]*pb.Stream) {
	//TODO(tcfw) merge streamDiff and nodeDiff
	added := map[int32]*pb.Stream{}
	deleted := map[int32]*pb.Stream{}
	same := map[int32]*pb.Stream{}

	streamMap := map[int32]*pb.Stream{}
	for _, sStream := range s.streams {
		streamMap[sStream.Id] = sStream
	}

	//Find new and same
	for _, nStream := range nStreams {
		if _, ok := streamMap[nStream.Id]; ok {
			same[nStream.Id] = nStream
		} else {
			added[nStream.Id] = nStream
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

func (s *basicScheduler) nodeDiff(nNodes []*pb.Node) (map[string]*pb.Node, map[string]*pb.Node) {
	added := map[string]*pb.Node{}
	deleted := map[string]*pb.Node{}
	same := map[string]*pb.Node{}

	//Find new and same
	for _, nNode := range nNodes {
		if _, ok := s.nodes[nNode.Id]; ok {
			same[nNode.Id] = nNode
		} else {
			added[nNode.Id] = nNode
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
		nodes:           map[string]*pb.Node{},
		streams:         []*pb.Stream{},
		bindings:        []*pb.Binding{},
		nf:              nf,
		sf:              sf,
		observeInternal: 5 * time.Second,
	}
}

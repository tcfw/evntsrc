package ttlscheduler

import (
	"errors"
	"reflect"
	"sync"
	"time"
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
	GetStreams() ([]int32, error)
}

//Scheduler basic structure of a scheduler
type Scheduler interface {
	NodeBindings(node) ([]*binding, error)
	BindStream(int32) (*binding, error)
	Optimise() error
	Observe() error
}

type basicScheduler struct {
	nodes    map[int]*node
	streams  []int32
	bindings []*binding
	nf       NodeFetcher
	sf       StreamFetcher
	once     bool

	lock sync.RWMutex
}

func (s *basicScheduler) NodeBindings(node node) ([]*binding, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	nodeBindings := []*binding{}

	for _, bind := range s.bindings {
		if bind.Node == nil {
			continue
		}
		if bind.Node.ID == node.ID {
			nodeBindings = append(nodeBindings, bind)
		}
	}

	return nodeBindings, nil
}

func (s *basicScheduler) GetNodes() map[int]*node {
	return s.nodes
}

//BindStream allocates a stream to a node worker in an "optimal" pattern
func (s *basicScheduler) BindStream(id int32) (*binding, error) {
	nCount := len(s.nodes)

	if nCount == 0 {
		return nil, errors.New("No nodes to schedule on")
	}
	if nCount == 1 {
		binding := &binding{Stream: &stream{id, 0}, Node: s.nodes[int(reflect.ValueOf(s.nodes).MapKeys()[0].Int())]}
		return binding, nil
	}

	nScores := map[int]int{}
	for i, node := range s.nodes {
		nScores[i] = NodeScore(s, *node)
	}

	lowestScore, lowestNode := -1, -1
	for node, score := range nScores {
		if lowestNode == -1 || score < lowestScore {
			lowestNode = node
			lowestScore = score
		}
	}

	binding := &binding{Stream: &stream{id, 0}, Node: s.nodes[lowestNode]}
	return binding, nil
}

//NodeScore calculates a load score for a particular node
func NodeScore(s Scheduler, node node) int {
	score := 0
	nodeBinds, _ := s.NodeBindings(node)

	for _, binding := range nodeBinds {
		score += 10 + binding.Stream.MsgRate
	}

	return score
}

//Optimise @TODO
func (s *basicScheduler) Optimise() error {
	return nil
}

func (s *basicScheduler) Observe() error {
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

		s.lock.Unlock()

		if s.once {
			return nil
		}

		time.Sleep(5 * time.Second)
	}
}

func (s *basicScheduler) observeNodes(nNodes []*node) {
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
			binding, _ := s.BindStream(resStream.ID)
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

func (s *basicScheduler) observeStreams(nStreams []int32) {
	addedStreams, deletedStreams := s.streamDiff(nStreams)
	for id := range addedStreams {
		binding, _ := s.BindStream(id)
		s.bindings = append(s.bindings, binding)
		s.streams = append(s.streams, id)
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
				if stream == id {
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
func (s *basicScheduler) streamDiff(nStreams []int32) (map[int32]int32, map[int32]int32) {
	added := map[int32]int32{}
	deleted := map[int32]int32{}
	same := map[int32]int32{}

	streamMap := map[int32]int32{}
	for _, stream := range s.streams {
		streamMap[stream] = stream
	}

	//Find new and same
	for _, nStream := range nStreams {
		if _, ok := streamMap[nStream]; ok {
			same[nStream] = nStream
		} else {
			added[nStream] = nStream
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
		nodes:    map[int]*node{},
		streams:  []int32{},
		bindings: []*binding{},
		nf:       nf,
		sf:       sf,
	}
}

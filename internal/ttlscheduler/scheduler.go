package ttlscheduler

import (
	"errors"
	"sync"
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
	GetStreams() ([]*int32, error)
}

//Scheduler basic structure of a scheduler
type Scheduler interface {
	NodeBindings(node) ([]*binding, error)
	BindStream(int32) (*binding, error)
	Optimise() error
	Observe() error
}

type basicScheduler struct {
	nodes    []*node
	streams  []int32
	bindings []*binding
	nf       NodeFetcher
	sf       StreamFetcher

	lock sync.RWMutex
}

func (s *basicScheduler) NodeBindings(node node) ([]*binding, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	nodeBindings := []*binding{}

	for _, bind := range s.bindings {
		if bind.Node.ID == node.ID {
			nodeBindings = append(nodeBindings, bind)
		}
	}

	return nodeBindings, nil
}

//BindStream allocates a stream to a node worker in an "optimal" pattern
func (s *basicScheduler) BindStream(id int32) (*binding, error) {
	nCount := len(s.nodes)

	if nCount == 0 {
		return nil, errors.New("No nodes to schedule on")
	}
	if nCount == 1 {
		binding := &binding{Stream: &stream{id, 0}, Node: s.nodes[0]}
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
	return nil
}

//NewScheduler constructs a new scheduler which can assign streams work nodes
func NewScheduler(nf NodeFetcher, sf StreamFetcher) Scheduler {
	return &basicScheduler{
		nodes:    []*node{},
		streams:  []int32{},
		bindings: []*binding{},
		nf:       nf,
		sf:       sf,
	}
}

package ttlscheduler

import (
	pb "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
)

type basicNodeFetcher struct{}

func (bnf *basicNodeFetcher) GetNodes() ([]*pb.Node, error) {
	return []*pb.Node{}, nil
}

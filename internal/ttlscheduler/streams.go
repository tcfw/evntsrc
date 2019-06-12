package ttlscheduler

import (
	pb "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
)

type basicStreamFetcher struct{}

func (bsf *basicStreamFetcher) GetStreams() ([]*pb.Stream, error) {
	return []*pb.Stream{}, nil
}

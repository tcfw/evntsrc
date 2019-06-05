package ttlscheduler

type basicNodeFetcher struct{}

func (bnf *basicNodeFetcher) GetNodes() ([]*node, error) {
	return []*node{}, nil
}

package ttlscheduler

type basicStreamFetcher struct{}

func (bsf *basicStreamFetcher) GetStreams() ([]*stream, error) {
	return []*stream{}, nil
}

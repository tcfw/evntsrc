package ttlscheduler

type basicStreamFetcher struct{}

func (bsf *basicStreamFetcher) GetStreams() ([]*int32, error) {
	return []*int32{}, nil
}

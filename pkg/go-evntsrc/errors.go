package evntsrc

import "time"

//pubError attempts to write the err to the client error channel
//gives up after 10 seconds
func (api *APIClient) pubError(err error) {
	go func() {
		select {
		case api.errCh <- err:
		case <-time.After(10 * time.Second):
		}
	}()
}

//Errors provides a channel for reading errors from
//API calls and subscription issues
func (api *APIClient) Errors() <-chan error {
	return api.errCh
}

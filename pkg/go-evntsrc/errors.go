package evntsrc

import "time"

//pubError attempts to write the err to the client error channel
//gives up after 10 seconds
func (api *APIClient) pubError(err error) {
	go func() {
		select {
		case api.Errors <- err:
		case <-time.After(10 * time.Second):
		}
	}()
}

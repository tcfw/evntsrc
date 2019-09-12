package websocks

import (
	"fmt"
	"sync"
	"time"
)

//Ack the basic ack struct
type Ack struct {
	ackType string
	at      time.Time
}

//Acks a group of acks with syncs
type Acks struct {
	acks sync.Map
	stop chan struct{}
	cond *sync.Cond
}

//NewAcks constructs a new map of acks with syncs
func NewAcks() *Acks {
	acks := &Acks{
		acks: sync.Map{},
		stop: make(chan struct{}),
		cond: sync.NewCond(&sync.Mutex{}),
	}
	go acks.GC()
	return acks
}

//GC ranges over the sync map to clear out acks older than 5 minutes
func (a *Acks) GC() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			a.acks.Range(func(k interface{}, v interface{}) bool {
				ack := v.(*Ack)
				if time.Now().Add(-5 * time.Minute).Before(ack.at) {
					a.acks.Delete(k)
				}
				return true
			})
			break
		case <-a.stop:
			ticker.Stop()
			return
		}
	}
}

//Has checks to see if a map key exists
func (a *Acks) Has(k string) bool {
	_, ok := a.acks.Load(k)
	return ok
}

//WaitForKey checks if a key exists otherwise waits for a change to
//the map and then checks again
func (a *Acks) WaitForKey(k string) {
	a.cond.L.Lock()
	defer a.cond.L.Unlock()

	for !a.Has(k) {
		a.cond.Wait()
	}
}

//WaitForKeyWithTimeout waits for the key to exist or times out after
//specified time
func (a *Acks) WaitForKeyWithTimeout(k string, d time.Duration) error {
	done := make(chan struct{})

	go func() {
		//TODO(tcfw) fix goroutine leak if key is never found
		a.WaitForKey(k)
		close(done)
	}()

	select {
	case <-time.After(d):
		return fmt.Errorf("Timed out after %s", d.String())
	case <-done:
		//change occured
	}

	return nil
}

//Add amends a key to the map and singles anything waiting for a
//key change
func (a *Acks) Add(key string, ackType string) *Ack {
	ac := &Ack{ackType: ackType, at: time.Now()}
	a.acks.Store(key, ac)
	a.cond.Broadcast()
	return ac
}

//Stop terminates the GC loop
func (a *Acks) Stop() {
	close(a.stop)
}

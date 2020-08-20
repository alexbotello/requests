package pool

import "sync"

// Requestor represents an object that will make an HTTP request
type Requestor interface {
	Request()
}

// RequestPool represents a pool of concurrent HTTP requests
//
// A sync.WaitGroup is used to keep RequestPool open until all workers have finished
// processing.
//
// A slice of Requestor objects will be passed into the PoolWorkers so they may handle
// the actual Request call
//
// A channel of type Requestor that will handle communicating the Requestors from
// the RequestPool --> poolWorker
type RequestPool struct {
	wg         *sync.WaitGroup
	requestors []Requestor
	channel    chan Requestor
}

// NewRequestPool instantiates a new RequestPool
func NewRequestPool(rqs []Requestor) *RequestPool {
	var w sync.WaitGroup
	ch := make(chan Requestor)
	return &RequestPool{wg: &w, channel: ch, requestors: rqs}
}

// Starts the RequestPool -- the pool will block until all poolWorkers
// have finished within the WaitGroup
func (rp RequestPool) Start() {
	rp.wg.Add(len(rp.requestors))
	go rp.SpawnWorkers()
	rp.wg.Wait()
}

// SpawnWorkers will spawn n number of poolWorkers
func (rp RequestPool) SpawnWorkers() {
	defer close(rp.channel)
	numofWorkers := len(rp.requestors)

	for i := 0; i < numofWorkers; i++ {
		go newPoolWorker(rp.wg, rp.channel)
	}
	for _, requestor := range rp.requestors {
		rp.channel <- requestor
	}
}

// poolWorker represents a concurrent task spawned by a RequestPool
type poolWorker struct{}

// poolWorker will Work until request is finished
func (w poolWorker) Work(wg *sync.WaitGroup, receiver chan Requestor) {
	for {
		requestor, ok := <-receiver
		if !ok {
			wg.Done()
			return
		}
		w.makeRequest(requestor)
	}
}

func (w poolWorker) makeRequest(r Requestor) {
	r.Request()
}

// instantiates a new poolWorker, this function should only be used by the
// RequestPool
func newPoolWorker(wg *sync.WaitGroup, receiver chan Requestor) {
	w := &poolWorker{}
	go w.Work(wg, receiver)
}

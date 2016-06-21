package godispatch

import (
	"sync"
)

// Dispatcher dispatches work requests
type Dispatcher struct {
	sync.RWMutex
	MasterWorkerMap map[Master]*Worker // There can be multiple Masters with the same MasterID but different SlaveID
	WorkHandler     Handler
}

// NewDispatcher returns a Dispatcher instance
func NewDispatcher(h interface{}) *Dispatcher {
	return &Dispatcher{
		MasterWorkerMap: make(map[Master]*Worker),
		WorkHandler:     h.(Handler),
	}
}

// Dispatch work to worker that matches WorkID
func (d *Dispatcher) Dispatch(w Work, m Master) {
	d.Lock()
	_, ok := d.MasterWorkerMap[m]

	if ok == false { // meter is not in meter map
		d.MasterWorkerMap[m] = NewWorker()
		d.MasterWorkerMap[m].Start(d)
	}
	workChannel := d.MasterWorkerMap[m].WorkChannel
	d.Unlock()

	// Send work to worker's Work Channel
	workChannel <- w
}

// Close closes all the workers' WorkChannels
func (d *Dispatcher) Close() {
	// Closing connections...
	d.RLock()
	for _, worker := range d.MasterWorkerMap {
		close(worker.WorkChannel) // dispatcher closes worker's channel
		worker.QuitChan <- true   // sends worker quit signal
	}
	d.RUnlock()
}

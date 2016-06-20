package godispatch

import (
	"sync"
)

// Master has its own ID and its Worker ID
type Master struct {
	MasterID, WorkerID string
}

// Dispatcher dispatches work requests
type Dispatcher struct {
	sync.RWMutex
	MasterWorkerMap map[Master]*Worker // There can be multiple Masters with the same MasterID but different WorkerID
	WorkHandler     interface{}
}

// Work received by Workers
type Work struct {
	MasterID  string `json:"MasterID"`
	SlaveID   string `json:"SlaveID"`
	EventTime int64  `json:"EventTime"`
}

// NewDispatcher returns a Dispatcher instance
func NewDispatcher(i interface{}) Dispatcher {
	return &Dispatcher{
		MasterWorkerMap: make(map[Master]*Worker),
		WorkHandler:     i,
	}
}

// Dispatch work to worker that matches WorkID
func (d *Dispatcher) Dispatch(w Work) {
	m := Master{w.MasterID, w.WorkerID}

	d.Lock()
	_, ok := d.MastMasterWorkerMaperMap[m]

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

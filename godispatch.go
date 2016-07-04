package godispatch

import (
	"flag"
	"io/ioutil"
	"log"
	"sync"
)

var debug = flag.Bool("d", false, "Set to True for debug info, default is false")

// Dispatcher dispatches work to worker
type Dispatcher struct {
	sync.RWMutex
	MasterWorkerMap map[Master]*Worker // There can be multiple Masters with the same MasterID but different SlaveID
	WorkHandler     Handler
}

// NewDispatcher returns a Dispatcher instance
func NewDispatcher(h interface{}) *Dispatcher {
	flag.Parse()
	if !*debug {
		log.SetOutput(ioutil.Discard)
	}
	return &Dispatcher{
		MasterWorkerMap: make(map[Master]*Worker),
		WorkHandler:     h.(Handler),
	}
}

// Dispatch work to worker that matches WorkID
func (d *Dispatcher) Dispatch(w Work, m Master) {

	d.Lock()
	_, ok := d.MasterWorkerMap[m]

	if ok == false { // Master is not in map
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
	log.Println("Closing Work Channels...")
	d.RLock()
	for _, worker := range d.MasterWorkerMap {
		close(worker.WorkChannel) // dispatcher closes worker's channel
		worker.QuitChan <- true   // sends worker quit signal
	}
	d.RUnlock()
	log.Println("All Work Channel(s) Closed")
}

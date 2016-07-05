package godispatch

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// Dispatcher dispatches work to worker
type Dispatcher struct {
	sync.RWMutex
	MasterWorkerMap map[Master]*Worker
	WorkHandler     Handler
	debug           bool // Set to true for godispatch debug info
}

// NewDispatcher returns a Dispatcher instance
func NewDispatcher(h interface{}, debugOption ...bool) *Dispatcher {
	d := &Dispatcher{
		MasterWorkerMap: make(map[Master]*Worker),
		WorkHandler:     h.(Handler),
		debug:           false, // Default is false
	}

	log.SetOutput(ioutil.Discard)

	if len(debugOption) > 0 {
		d.debug = debugOption[0]
	}

	if d.debug == true {
		log.SetOutput(os.Stdout)
	}

	log.Println("Dispatcher Created")

	return d
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

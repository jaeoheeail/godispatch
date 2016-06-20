package godispatch

import (
	"sync"
)

// Worker has a work channel to receive work
type Worker struct {
	WorkChannel chan Work
	QuitChan    chan bool
}

// NewWorker returns a Worker instance
func NewWorker() *Worker {
	return &Worker{
		WorkChannel: make(chan Work),
		QuitChan:    make(chan bool),
	}
}

// Start receiving work from work channel and begin processing
func (w *Worker) Start(d *Dispatcher) {
	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case work := <-w.WorkChannel:
				// Worker: Received work
				wg.Add(1)
				go func() {
					defer wg.Done()
					d.WorkHandler.Handle(work)
				}()
				wg.Wait()
			case <-w.QuitChan:
				// Worker: Received quit signal
				return
			}
		}
	}()
}

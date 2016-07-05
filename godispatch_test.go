package godispatch

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var handleLock *sync.RWMutex

// MyWork is a Work struct used for testing
type MyWork struct {
	MasterID string
	WorkID   string
	Done     bool
}

// MyMaster is a Master struct used for testing
type MyMaster struct {
	MasterID string
}

// MyHandler handles Work for testing
type MyHandler interface {
	Handle(e Work)
}

// MyHandlerStruct is a Handler struct used for testing
type MyHandlerStruct struct {
}

// Handle handles work for testing
func (h *MyHandlerStruct) Handle(g Work) {
	w, ok := g.(MyWork) // Type Assertion
	if !ok {
		return
	}

	// Do some work...
	w.Done = true

	handleLock.Lock()
	finishedWork = append(finishedWork, w)
	handleLock.Unlock()
}

// MakeHandler returns MyHandler interface
func MakeHandler() MyHandler {
	handleLock = new(sync.RWMutex)
	return &MyHandlerStruct{}
}

var _ Work = (*MyWork)(nil)
var _ Master = (*MyMaster)(nil)
var _ Handler = (*MyHandlerStruct)(nil)

var finishedWork []MyWork

func TestDispatcherNoDebugging(t *testing.T) {
	h := MakeHandler()
	d := NewDispatcher(h)

	masters := []MyMaster{
		{MasterID: "1"},
		{MasterID: "2"},
		{MasterID: "3"},
		{MasterID: "4"},
	}

	for _, m := range masters {
		for i := 1; i < 5; i++ {
			w := MyWork{MasterID: m.MasterID, WorkID: strconv.Itoa(i), Done: false}
			go d.Dispatch(w, m)
		}
	}

	time.Sleep(10000000) // Wait until all work has been dispatched

	d.Lock()
	// Check MasterWorkerMap has 4 Masters
	assert.Equal(t, len(d.MasterWorkerMap), 4)
	d.Unlock()

	handleLock.Lock()
	for _, w := range finishedWork { // Check that all work is done
		assert.Equal(t, w.Done, true)
	}
	handleLock.Unlock()

	d.Close()

	// Check Workers' WorkChannel in MasterWorkerMap are closed
	for _, m := range masters {
		workChannel := d.MasterWorkerMap[m].WorkChannel
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic did not occur")
			}
		}()
		workChannel <- MyWork{} // Expect panic to occur
	}
}

func TestDispatcherWithDebugging(t *testing.T) {
	h := MakeHandler()
	d := NewDispatcher(h, true)

	masters := []MyMaster{
		{MasterID: "1"},
		{MasterID: "2"},
		{MasterID: "3"},
		{MasterID: "4"},
	}

	for _, m := range masters {
		for i := 1; i < 5; i++ {
			w := MyWork{MasterID: m.MasterID, WorkID: strconv.Itoa(i), Done: false}
			go d.Dispatch(w, m)
		}
	}

	time.Sleep(10000000) // Wait until all work has been dispatched

	d.Lock()
	// Check MasterWorkerMap has 4 Masters
	assert.Equal(t, len(d.MasterWorkerMap), 4)
	d.Unlock()

	handleLock.Lock()
	for _, w := range finishedWork { // Check that all work is done
		assert.Equal(t, w.Done, true)
	}
	handleLock.Unlock()

	d.Close()

	// Check Workers' WorkChannel in MasterWorkerMap are closed
	for _, m := range masters {
		workChannel := d.MasterWorkerMap[m].WorkChannel
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic did not occur")
			}
		}()
		workChannel <- MyWork{} // Expect panic to occur
	}
}

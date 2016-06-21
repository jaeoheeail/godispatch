package godispatch

import (
	"strconv"
	"testing"
	"time"

	//log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

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
	time.Sleep(5000)
	w.Done = true

	finishedWork = append(finishedWork, w)
}

// MakeHandler returns MyHandler interface
func MakeHandler() MyHandler {
	return &MyHandlerStruct{}
}

var _ Work = (*MyWork)(nil)
var _ Master = (*MyMaster)(nil)
var _ Handler = (*MyHandlerStruct)(nil)

var finishedWork []MyWork

func TestDispatcher(t *testing.T) {
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

	time.Sleep(1000000)

	// Check MasterWorkerMap has 4 Masters
	assert.Equal(t, len(d.MasterWorkerMap), 4)

	for _, w := range finishedWork { // Check that all work is done
		assert.Equal(t, w.Done, true)
	}

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

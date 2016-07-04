[![](https://img.shields.io/badge/coverage-100%25-brightgreen.svg)]()

# godispatch
> Simple golang package that dispatches work to workers 
> Inspired by [http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)

## Installation
`go get github.com/jaeoheeail/godispatch`

## Features
* Master Worker Map - Each Worker (Value) is tagged to a Master (Key) as a Key-Value pair in a map
* Provides sequential processing 

## Example
```golang
package main

import (
	"github.com/jaeoheeail/godispatch"
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

func main() {
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

	d.Close()
}
```

## Logging
* Debug log is turned **off** by default
* To turn on, add the following flag: `-d=true`

## Testing
* Refer to [`godispatch_test.go`](https://github.com/jaeoheeail/godispatch/blob/master/godispatch_test.go) for more details.
* Run `go test`
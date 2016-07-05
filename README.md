![alt text](https://img.shields.io/badge/coverage-100%25-brightgreen.svg)

# godispatch
> Simple golang package that dispatches work to workers 
> Inspired by [http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)

## Installation
`go get github.com/jaeoheeail/godispatch`

## Features
* Master Worker Map - Each Worker (Value) is tagged to a Master (Key) as a Key-Value pair in a map
* Provides sequential processing 

## Logging
* Debug log is turned **off** by default
* Dispatcher has a private attribute `debug` that is set to false by default
	```
	type Dispatcher struct {
		sync.RWMutex
		MasterWorkerMap map[Master]*Worker
		WorkHandler     Handler
		debug           bool // Set to true for godispatch debug info
	}
	```

* When initializing Dispatcher with `func NewDispatcher(h interface{}, debugOption ...bool) *Dispatcher`, use `NewDispatcher(h, true)` for debug logs
		

## Example

1. Define the following structs for Work and Master
```
type MyWork struct {
	MasterID string
	WorkID   string
	Done     bool
}

type MyMaster struct {
	MasterID string
}
```

2. Define interface for Handler and include `Handle` method and a constructor method (e.g. `MakeHandler`)

```
type MyHandler interface {
	Handle(e Work)
}

type MyHandlerStruct struct {
}

func (h *MyHandlerStruct) Handle(g Work) {
	w, ok := g.(MyWork) // Type Assertion
	if !ok {
		return
	}

	// Do some work...
	w.Done = true

}

// MakeHandler returns MyHandler interface
func MakeHandler() MyHandler {
	handleLock = new(sync.RWMutex)
	return &MyHandlerStruct{}
}
```

3. Suppose you have 4 Masters, each with 4 Work to complete
```
var _ Work = (*MyWork)(nil)
var _ Master = (*MyMaster)(nil)
var _ Handler = (*MyHandlerStruct)(nil)

func main() {
	h := MakeHandler()
	d := NewDispatcher(h) // use d := NewDispatcher(h, true) for debug info

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

## Testing
* Refer to [`godispatch_test.go`](https://github.com/jaeoheeail/godispatch/blob/master/godispatch_test.go) for more details.
* Run `go test`
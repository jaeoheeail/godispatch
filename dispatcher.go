package godispatch

// Dispatcher dispatches work to workers
type Dispatcher interface {
	Dispatch(w Work)
	Close()
}

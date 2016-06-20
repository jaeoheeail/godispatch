package godispatch

// GoDispatcher dispatches work to workers
type GoDispatcher interface {
	Dispatch(w Work)
	Close()
}

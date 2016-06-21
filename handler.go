package godispatch

// Handler handles Work
type Handler interface {
	Handle(w Work)
}

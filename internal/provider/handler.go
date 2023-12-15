package provider

type Handler[T any] struct {
	Err  func(error)
	Done func()
	Resp func(T)
}

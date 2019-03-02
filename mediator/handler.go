package mediator

type Handler interface {
	Handle(interface{}) (interface{}, error)
}

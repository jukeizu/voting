package server

type Mediator interface {
	Send(interface{}) (interface{}, error)
}

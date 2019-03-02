package registration

type Mediator interface {
	Send(interface{}) (interface{}, error)
}

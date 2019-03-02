package mediator

import (
	"fmt"

	"github.com/jukeizu/voting/registration"
)

type Mediator struct {
	registerVoterCommandHandler registration.RegisterVoterCommandHandler
}

func New(registerVoterCommandHandler registration.RegisterVoterCommandHandler) Mediator {
	return Mediator{registerVoterCommandHandler}
}

func (m Mediator) Send(req interface{}) (interface{}, error) {
	switch t := req.(type) {
	case registration.RegisterVoterRequest:
		return m.registerVoterCommandHandler.Handle(req.(registration.RegisterVoterRequest))
	default:
		return nil, fmt.Errorf("mediator: unsupported request type: %T", t)
	}
}

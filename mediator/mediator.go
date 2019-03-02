package mediator

import (
	"fmt"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
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
	case *registrationpb.RegisterVoterRequest:
		return m.registerVoterCommandHandler.Handle(req.(*registrationpb.RegisterVoterRequest))
	default:
		return nil, fmt.Errorf("mediator: unsupported request type: %T", t)
	}
}

package mediator

import (
	"fmt"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
)

type Mediator struct {
	registerVoterCommandHandler Handler
}

func New(registerVoterCommandHandler Handler) Mediator {
	return Mediator{registerVoterCommandHandler}
}

func (m Mediator) Send(req interface{}) (interface{}, error) {
	switch t := req.(type) {
	case *registrationpb.RegisterVoterRequest:
		return m.registerVoterCommandHandler.Handle(req)
	default:
		return nil, fmt.Errorf("mediator: unsupported request type: %T", t)
	}
}

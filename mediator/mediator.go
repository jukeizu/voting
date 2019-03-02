package mediator

import (
	"fmt"

	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/jukeizu/voting/polls"
	"github.com/jukeizu/voting/registration"
)

type Mediator struct {
	registerVoterCommandHandler registration.RegisterVoterCommandHandler
	createPollCommandHandler    polls.CreatePollCommandHandler
}

func New(
	registerVoterCommandHandler registration.RegisterVoterCommandHandler,
	createPollCommandHandler polls.CreatePollCommandHandler,
) Mediator {
	return Mediator{
		registerVoterCommandHandler,
		createPollCommandHandler,
	}
}

func (m Mediator) Send(req interface{}) (interface{}, error) {
	switch t := req.(type) {
	case *registrationpb.RegisterVoterRequest:
		return m.registerVoterCommandHandler.Handle(req.(*registrationpb.RegisterVoterRequest))
	case *pollpb.CreatePollRequest:
		return m.createPollCommandHandler.Handle(req.(*pollpb.CreatePollRequest))
	default:
		return nil, fmt.Errorf("mediator: unsupported request type: %T", t)
	}
}

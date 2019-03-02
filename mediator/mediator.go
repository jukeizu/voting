package mediator

import (
	"fmt"

	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
)

type Mediator struct {
	registerVoterCommandHandler Handler
	createPollCommandHandler    Handler
	pollQueryHandler            Handler
}

func New(
	registerVoterCommandHandler Handler,
	createPollCommandHandler Handler,
	pollQueryHandler Handler,
) Mediator {
	return Mediator{
		registerVoterCommandHandler,
		createPollCommandHandler,
		pollQueryHandler,
	}
}

func (m Mediator) Send(req interface{}) (interface{}, error) {
	switch t := req.(type) {
	case *registrationpb.RegisterVoterRequest:
		return m.registerVoterCommandHandler.Handle(req)
	case *pollpb.CreatePollRequest:
		return m.createPollCommandHandler.Handle(req)
	case *pollpb.PollRequest:
		return m.pollQueryHandler.Handle(req)
	default:
		return nil, fmt.Errorf("mediator: unsupported request type: %T", t)
	}
}

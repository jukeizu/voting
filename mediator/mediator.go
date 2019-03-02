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
	optionsQueryHandler         Handler
	endPollCommandHandler       Handler
}

func New(
	registerVoterCommandHandler Handler,
	createPollCommandHandler Handler,
	pollQueryHandler Handler,
	optionsQueryHandler Handler,
	endPollCommandHandler Handler,
) Mediator {
	return Mediator{
		registerVoterCommandHandler,
		createPollCommandHandler,
		pollQueryHandler,
		optionsQueryHandler,
		endPollCommandHandler,
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
	case *pollpb.OptionsRequest:
		return m.optionsQueryHandler.Handle(req)
	case *pollpb.EndPollRequest:
		return m.endPollCommandHandler.Handle(req)
	default:
		return nil, fmt.Errorf("mediator: unsupported request type: %T", t)
	}
}

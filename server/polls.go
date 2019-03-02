package server

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
)

type PollsServer struct {
	mediator Mediator
}

func NewPollsServer(mediator Mediator) PollsServer {
	return PollsServer{mediator}
}

func (s PollsServer) Create(ctx context.Context, req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
	resp, err := s.mediator.Send(req)
	if err != nil {
		return nil, err
	}

	reply, ok := resp.(*pollpb.CreatePollReply)
	if !ok {
		return nil, nil
	}

	return reply, nil
}

func (s PollsServer) Poll(ctx context.Context, req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	resp, err := s.mediator.Send(req)
	if err != nil {
		return nil, err
	}

	reply, ok := resp.(*pollpb.PollReply)
	if !ok {
		return nil, nil
	}

	return reply, nil
}

func (s PollsServer) Options(ctx context.Context, req *pollpb.OptionsRequest) (*pollpb.OptionsReply, error) {
	resp, err := s.mediator.Send(req)
	if err != nil {
		return nil, err
	}

	reply, ok := resp.(*pollpb.OptionsReply)
	if !ok {
		return nil, nil
	}

	return reply, nil
}

func (s PollsServer) End(ctx context.Context, req *pollpb.EndPollRequest) (*pollpb.EndPollReply, error) {
	resp, err := s.mediator.Send(req)
	if err != nil {
		return nil, err
	}

	reply, ok := resp.(*pollpb.EndPollReply)
	if !ok {
		return nil, nil
	}

	return reply, nil
}

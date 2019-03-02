package main

import (
	"context"
	"net"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/jukeizu/voting/domain/entities"
	"github.com/jukeizu/voting/mediator"
	"github.com/jukeizu/voting/registration"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Server struct {
	logger     zerolog.Logger
	grpcServer *grpc.Server
	mediator   mediator.Mediator
}

func NewServer(logger zerolog.Logger, grpcServer *grpc.Server, mediator mediator.Mediator) Server {
	return Server{logger, grpcServer, mediator}
}

func (s Server) RegisterVoter(ctx context.Context, req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
	request := registration.RegisterVoterRequest{
		ExternalId: req.ExternalId,
		Username:   req.Username,
	}

	voter, err := s.voterRequest(request)
	if err != nil {
		return nil, err
	}

	return &registrationpb.RegisterVoterReply{Voter: voter}, nil
}

func (s Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.logger.Info().
		Str("transport", "grpc").
		Str("addr", addr).
		Msg("listening")

	return s.grpcServer.Serve(listener)
}

func (s Server) Stop() {
	if s.grpcServer == nil {
		return
	}

	s.logger.Info().
		Str("transport", "grpc").
		Msg("stopping")

	s.grpcServer.GracefulStop()
}

func (s Server) voterRequest(req interface{}) (*registrationpb.Voter, error) {
	resp, err := s.mediator.Send(req)
	if err != nil {
		return nil, err
	}

	voter, ok := resp.(*entities.Voter)
	if !ok {
		return nil, nil
	}

	replyVoter := registrationpb.Voter{
		Id:         voter.Id,
		ExternalId: voter.ExternalId,
		Username:   voter.Username,
		CanVote:    voter.CanVote,
	}

	return &replyVoter, nil
}

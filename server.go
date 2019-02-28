package main

import (
	"context"
	"net"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/jukeizu/voting/persistence"
	"github.com/jukeizu/voting/registration"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Server struct {
	logger     zerolog.Logger
	grpcServer *grpc.Server
	repository persistence.Repository
}

func NewServer(logger zerolog.Logger, grpcServer *grpc.Server, repository persistence.Repository) Server {
	return Server{logger, grpcServer, repository}
}

func (s Server) RegisterVoter(ctx context.Context, req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
	request := registration.RegisterVoterRequest{
		ExternalId: req.ExternalId,
		Username:   req.Username,
	}

	handler := registration.NewRegisterVoterCommandHandler(s.logger, s.repository)
	voter, err := handler.Handle(request)
	if err != nil {
		return nil, err
	}

	replyVoter := registrationpb.Voter{
		Id:         voter.Id,
		ExternalId: voter.ExternalId,
		Username:   voter.Username,
		CanVote:    voter.CanVote,
	}

	return &registrationpb.RegisterVoterReply{Voter: &replyVoter}, nil
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

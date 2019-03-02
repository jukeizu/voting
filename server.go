package main

import (
	"context"
	"net"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/jukeizu/voting/mediator"
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
	resp, err := s.mediator.Send(req)
	if err != nil {
		return nil, err
	}

	reply, ok := resp.(*registrationpb.RegisterVoterReply)
	if !ok {
		return nil, nil
	}

	return reply, nil
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

package session

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/sessionpb"
)

type GrpcServer struct {
	service Service
}

func NewGrpcServer(service Service) GrpcServer {
	return GrpcServer{service}
}

func (s GrpcServer) CurrentPoll(ctx context.Context, req *sessionpb.CurrentPollRequest) (*sessionpb.CurrentPollReply, error) {
	pollId, err := s.service.CurrentPoll(req.ServerId)
	if err != nil {
		return nil, err
	}

	return &sessionpb.CurrentPollReply{ServerId: req.ServerId, PollId: pollId}, nil
}

func (s GrpcServer) SetCurrentPoll(ctx context.Context, req *sessionpb.SetCurrentPollRequest) (*sessionpb.SetCurrentPollReply, error) {
	err := s.service.SetCurrentPoll(req.ServerId, req.PollId)
	if err != nil {
		return nil, err
	}

	return &sessionpb.SetCurrentPollReply{ServerId: req.ServerId, PollId: req.PollId}, nil
}

func (s GrpcServer) Options(ctx context.Context, req *sessionpb.OptionsRequest) (*sessionpb.OptionsReply, error) {
	return nil, nil
}

func (s GrpcServer) Ballot(ctx context.Context, req *sessionpb.BallotRequest) (*sessionpb.BallotReply, error) {
	return nil, nil
}

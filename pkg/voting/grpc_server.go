package voting

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
)

var _ votingpb.VotingServer = &GrpcServer{}

type GrpcServer struct {
	service Service
}

func NewGrpcServer(service Service) GrpcServer {
	return GrpcServer{service}
}

func (s GrpcServer) CreatePoll(ctx context.Context, req *votingpb.CreatePollRequest) (*votingpb.CreatePollReply, error) {
	pollReq := createPollRequestToPoll(req)

	poll, err := s.service.CreatePoll(pollReq)
	if err != nil {
		return nil, err
	}

	pollReply := toPbPoll(poll)

	return &votingpb.CreatePollReply{Poll: pollReply}, nil
}

func (s GrpcServer) Poll(ctx context.Context, req *votingpb.PollRequest) (*votingpb.PollReply, error) {
	poll, err := s.service.Poll(req.ShortId, req.ServerId)
	if err != nil {
		return nil, err
	}

	pollReply := toPbPoll(poll)

	return &votingpb.PollReply{Poll: pollReply}, nil
}

func (s GrpcServer) EndPoll(ctx context.Context, req *votingpb.EndPollRequest) (*votingpb.EndPollReply, error) {
	poll, err := s.service.EndPoll(req.ShortId, req.ServerId, req.RequesterId)
	if err != nil {
		return nil, err
	}

	pollReply := toPbPoll(poll)

	return &votingpb.EndPollReply{Poll: pollReply}, nil
}

func (s GrpcServer) Status(ctx context.Context, req *votingpb.StatusRequest) (*votingpb.StatusReply, error) {
	status, err := s.service.Status(req.ShortId, req.ServerId)
	if err != nil {
		return nil, err
	}

	reply := &votingpb.StatusReply{
		Poll:   toPbPoll(status.Poll),
		Voters: toPbVoters(status.Voters),
	}

	return reply, nil
}

func (s GrpcServer) Vote(ctx context.Context, req *votingpb.VoteRequest) (*votingpb.VoteReply, error) {
	return nil, nil
}

func (s GrpcServer) Count(ctx context.Context, req *votingpb.CountRequest) (*votingpb.CountReply, error) {
	return nil, nil
}

func createPollRequestToPoll(req *votingpb.CreatePollRequest) Poll {
	poll := Poll{
		ServerId:           req.ServerId,
		CreatorId:          req.CreatorId,
		Title:              req.Title,
		AllowedUniqueVotes: req.AllowedUniqueVotes,
	}

	for _, createOption := range req.Options {
		option := Option{
			Content: createOption.Content,
			Url:     createOption.Url,
		}

		poll.Options = append(poll.Options, option)
	}

	return poll
}

func toPbPoll(poll Poll) *votingpb.Poll {
	pbPoll := &votingpb.Poll{
		Id:                 poll.Id,
		ShortId:            poll.ShortId,
		ServerId:           poll.ServerId,
		CreatorId:          poll.CreatorId,
		Title:              poll.Title,
		AllowedUniqueVotes: poll.AllowedUniqueVotes,
		HasEnded:           poll.HasEnded,
		Options:            toPbOptions(poll.Options),
	}

	return pbPoll
}

func toPbOptions(options []Option) []*votingpb.Option {
	pbOptions := []*votingpb.Option{}

	for _, option := range options {
		pbOption := &votingpb.Option{
			Id:      option.Id,
			Content: option.Content,
			Url:     option.Url,
		}

		pbOptions = append(pbOptions, pbOption)
	}

	return pbOptions
}

func toPbVoters(voters []Voter) []*votingpb.Voter {
	pbVoters := []*votingpb.Voter{}

	for _, voter := range voters {
		pbVoter := &votingpb.Voter{
			Id:       voter.Id,
			Username: voter.Username,
		}

		pbVoters = append(pbVoters, pbVoter)
	}

	return pbVoters
}

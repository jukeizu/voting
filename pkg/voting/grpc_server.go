package voting

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, toStatusErr(err)
	}

	pollReply := toPbPoll(poll)

	return &votingpb.CreatePollReply{Poll: pollReply}, nil
}

func (s GrpcServer) Poll(ctx context.Context, req *votingpb.PollRequest) (*votingpb.PollReply, error) {
	poll, err := s.service.Poll(req.ShortId, req.VoterId, req.ServerId)
	if err != nil {
		return nil, toStatusErr(err)
	}

	pollReply := toPbPoll(poll)

	return &votingpb.PollReply{Poll: pollReply}, nil
}

func (s GrpcServer) EndPoll(ctx context.Context, req *votingpb.EndPollRequest) (*votingpb.EndPollReply, error) {
	poll, err := s.service.EndPoll(req.ShortId, req.ServerId, req.RequesterId)
	if err != nil {
		return nil, toStatusErr(err)
	}

	pollReply := toPbPoll(poll)

	return &votingpb.EndPollReply{Poll: pollReply}, nil
}

func (s GrpcServer) Status(ctx context.Context, req *votingpb.StatusRequest) (*votingpb.StatusReply, error) {
	status, err := s.service.Status(req.ShortId, req.ServerId)
	if err != nil {
		return nil, toStatusErr(err)
	}

	reply := &votingpb.StatusReply{
		Poll:       toPbPoll(status.Poll),
		VoterCount: status.VoterCount,
	}

	return reply, nil

}

func (s GrpcServer) Voters(req *votingpb.VotersRequest, stream votingpb.Voting_VotersServer) error {
	voters, err := s.service.Voters(req.ShortId, req.ServerId)
	if err != nil {
		return toStatusErr(err)
	}

	for _, voter := range voters {
		err := stream.Send(toPbVoter(voter))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s GrpcServer) Vote(ctx context.Context, req *votingpb.VoteRequest) (*votingpb.VoteReply, error) {
	voteRequest := VoteRequest{
		ShortId:  req.ShortId,
		ServerId: req.ServerId,
		Voter:    toVoter(req.Voter),
		Options:  toBallotOptions(req.Options),
	}

	voteReply, err := s.service.Vote(voteRequest)
	if err != nil {
		return nil, toStatusErr(err)
	}

	return toPbVoteReply(voteReply), nil
}

func (s GrpcServer) Count(ctx context.Context, req *votingpb.CountRequest) (*votingpb.CountReply, error) {
	countRequest := CountRequest{
		ShortId:    req.ShortId,
		ServerId:   req.ServerId,
		NumToElect: int(req.NumToElect),
		Method:     req.Method,
		ToExclude:  req.ToExclude,
	}

	countResult, err := s.service.Count(countRequest)
	if err != nil {
		return nil, toStatusErr(err)
	}

	return toPbCountReply(countResult), nil
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

func toBallotOptions(pbOptions []*votingpb.BallotOption) []BallotOption {
	ballotOptions := []BallotOption{}

	for _, pbOption := range pbOptions {
		ballotOption := BallotOption{
			Rank:     pbOption.Rank,
			OptionId: pbOption.OptionId,
		}

		ballotOptions = append(ballotOptions, ballotOption)
	}

	return ballotOptions
}

func toVoter(pbVoter *votingpb.Voter) Voter {
	if pbVoter == nil {
		return Voter{}
	}

	voter := Voter{
		ExternalId: pbVoter.Id,
		Username:   pbVoter.Username,
	}

	return voter
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
		pbOptions = append(pbOptions, toPbOption(option))
	}

	return pbOptions
}

func toPbOption(option Option) *votingpb.Option {
	pbOption := &votingpb.Option{
		Id:      option.Id,
		Content: option.Content,
		Url:     option.Url,
	}

	return pbOption
}

func toPbVoteReply(voteReply VoteReply) *votingpb.VoteReply {
	pbVoteReply := &votingpb.VoteReply{
		Success: voteReply.Success,
		Message: voteReply.Message,
	}

	for _, voteReplyOption := range voteReply.Options {
		pbVoteReply.Options = append(pbVoteReply.Options, toPbVoteReplyOption(voteReplyOption))
	}

	return pbVoteReply
}

func toPbVoteReplyOption(voteReplyOption VoteReplyOption) *votingpb.VoteReplyOption {
	pbVoteReplyOption := &votingpb.VoteReplyOption{
		Rank:   voteReplyOption.Rank,
		Option: toPbOption(voteReplyOption.Option),
	}

	return pbVoteReplyOption
}

func toPbVoter(voter Voter) *votingpb.Voter {
	pbVoter := &votingpb.Voter{
		Id:       voter.Id,
		Username: voter.Username,
	}

	return pbVoter
}

func toPbCountReply(countResult CountResult) *votingpb.CountReply {
	countReply := &votingpb.CountReply{
		Success:   countResult.Success,
		Message:   countResult.Message,
		Poll:      toPbPoll(countResult.Poll),
		Method:    countResult.Method,
		Events:    toPbCountEvents(countResult.Events),
		Summaries: toPbCountEvents(countResult.Summaries),
	}

	for _, elected := range countResult.Elected {
		countReply.Elected = append(countReply.Elected, toPbVoteReplyOption(elected))
	}

	return countReply
}

func toPbCountEvents(countEvents []CountEvent) []*votingpb.CountEvent {
	pbCountEvents := []*votingpb.CountEvent{}

	for _, countEvent := range countEvents {
		pbCountEvent := &votingpb.CountEvent{
			Description: countEvent.Description,
		}

		pbCountEvents = append(pbCountEvents, pbCountEvent)
	}

	return pbCountEvents
}

func toStatusErr(err error) error {
	switch err.(type) {
	case ValidationError:
		return status.Error(codes.InvalidArgument, err.Error())
	case NotFoundError:
		return status.Error(codes.NotFound, err.Error())
	}
	return err
}

package voting

import (
	"context"
	"time"

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

func (s GrpcServer) VoterPoll(ctx context.Context, req *votingpb.VoterPollRequest) (*votingpb.PollReply, error) {
	poll, err := s.service.VoterPoll(req.VoterId, req.ServerId)
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

func (s GrpcServer) OpenPoll(ctx context.Context, req *votingpb.OpenPollRequest) (*votingpb.OpenPollReply, error) {
	openPollResult, err := s.service.OpenPoll(req.ShortId, req.ServerId, req.RequesterId, time.Unix(req.Expires, 0))
	if err != nil {
		return nil, toStatusErr(err)
	}

	pollReply := toPbPoll(openPollResult.Poll)

	return &votingpb.OpenPollReply{
		Poll:               pollReply,
		PreviouslyEnded:    openPollResult.PreviouslyEnded,
		PreviousExpiration: openPollResult.PreviousExpiration.UTC().Unix(),
	}, nil
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

func (s GrpcServer) Export(ctx context.Context, req *votingpb.ExportRequest) (*votingpb.ExportReply, error) {
	exportRequest := ExportRequest{
		ShortId:    req.ShortId,
		ServerId:   req.ServerId,
		NumToElect: int(req.NumToElect),
		Method:     req.Method,
		ToExclude:  req.ToExclude,
	}

	result, err := s.service.Export(exportRequest)
	if err != nil {
		return nil, toStatusErr(err)
	}
	return toPbExportReply(result), nil
}

func createPollRequestToPoll(req *votingpb.CreatePollRequest) Poll {
	poll := Poll{
		ServerId:           req.ServerId,
		CreatorId:          req.CreatorId,
		Title:              req.Title,
		AllowedUniqueVotes: req.AllowedUniqueVotes,
		Expires:            time.Unix(req.Expires, 0),
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
		ExternalId: pbVoter.ExternalId,
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
		HasEnded:           poll.HasEnded(),
		Options:            toPbOptions(poll.Options),
		Expires:            poll.Expires.UTC().Unix(),
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
		Id:         voter.Id,
		Username:   voter.Username,
		ExternalId: voter.ExternalId,
	}

	return pbVoter
}

func toPbCountReply(countResult CountResult) *votingpb.CountReply {
	countReply := &votingpb.CountReply{
		Success: countResult.Success,
		Message: countResult.Message,
		Poll:    toPbPoll(countResult.Poll),
		Method:  countResult.Method,
		Events:  toPbCountEvents(countResult.Events),
		Rounds:  toPbRoundSummaries(countResult.Rounds),
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
			Type:        countEvent.Type,
			Description: countEvent.Description,
		}

		pbCountEvents = append(pbCountEvents, pbCountEvent)
	}

	return pbCountEvents
}

func toPbRoundSummaries(rounds []RoundSummary) []*votingpb.RoundSummary {
	roundSummaries := []*votingpb.RoundSummary{}

	for _, r := range rounds {
		roundSummary := &votingpb.RoundSummary{
			Number:     int32(r.Number),
			Excess:     r.Excess,
			Surplus:    r.Surplus,
			Quota:      r.Quota,
			Candidates: toPbCandidateSummaries(r.Candidates),
		}

		roundSummaries = append(roundSummaries, roundSummary)
	}

	return roundSummaries
}

func toPbCandidateSummaries(candidates []CandidateSummary) []*votingpb.CandidateSummary {
	candidateSummaries := []*votingpb.CandidateSummary{}

	for _, c := range candidates {
		candidateSummary := &votingpb.CandidateSummary{
			Id:     c.Id,
			Name:   c.Name,
			Rank:   int32(c.Rank),
			Votes:  c.Votes,
			Weight: c.Weight,
			Status: c.Status,
		}

		candidateSummaries = append(candidateSummaries, candidateSummary)
	}

	return candidateSummaries
}

func toPbExportReply(exportResult ExportResult) *votingpb.ExportReply {
	exportReply := &votingpb.ExportReply{
		Content: exportResult.Content,
		Poll:    toPbPoll(exportResult.Poll),
		Method:  exportResult.Method,
	}

	return exportReply
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

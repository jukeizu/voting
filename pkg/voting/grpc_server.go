package voting

import (
	"context"
	"time"

	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ votingpb.VotingServiceServer = &GrpcServer{}

type GrpcServer struct {
	service Service
}

func NewGrpcServer(service Service) GrpcServer {
	return GrpcServer{service}
}

func (s GrpcServer) CreatePoll(ctx context.Context, req *votingpb.CreatePollRequest) (*votingpb.CreatePollResponse, error) {
	pollReq := createPollRequestToPoll(req)

	poll, err := s.service.CreatePoll(pollReq)
	if err != nil {
		return nil, toStatusErr(err)
	}

	pollResponse := toPbPoll(poll)

	return &votingpb.CreatePollResponse{Poll: pollResponse}, nil
}

func (s GrpcServer) Poll(ctx context.Context, req *votingpb.PollRequest) (*votingpb.PollResponse, error) {
	poll, err := s.service.Poll(req.ShortId, req.VoterId, req.ServerId)
	if err != nil {
		return nil, toStatusErr(err)
	}

	pollResponse := toPbPoll(poll)

	return &votingpb.PollResponse{Poll: pollResponse}, nil
}

func (s GrpcServer) VoterPoll(ctx context.Context, req *votingpb.VoterPollRequest) (*votingpb.VoterPollResponse, error) {
	poll, err := s.service.VoterPoll(req.VoterId, req.ServerId)
	if err != nil {
		return nil, toStatusErr(err)
	}

	pollResponse := toPbPoll(poll)

	return &votingpb.VoterPollResponse{Poll: pollResponse}, nil
}

func (s GrpcServer) EndPoll(ctx context.Context, req *votingpb.EndPollRequest) (*votingpb.EndPollResponse, error) {
	poll, err := s.service.EndPoll(req.ShortId, req.ServerId, req.RequesterId)
	if err != nil {
		return nil, toStatusErr(err)
	}

	pollResponse := toPbPoll(poll)

	return &votingpb.EndPollResponse{Poll: pollResponse}, nil
}

func (s GrpcServer) OpenPoll(ctx context.Context, req *votingpb.OpenPollRequest) (*votingpb.OpenPollResponse, error) {
	openPollResult, err := s.service.OpenPoll(req.ShortId, req.ServerId, req.RequesterId, time.Unix(req.Expires, 0))
	if err != nil {
		return nil, toStatusErr(err)
	}

	pollResponse := toPbPoll(openPollResult.Poll)

	return &votingpb.OpenPollResponse{
		Poll:               pollResponse,
		PreviouslyEnded:    openPollResult.PreviouslyEnded,
		PreviousExpiration: openPollResult.PreviousExpiration.UTC().Unix(),
	}, nil
}

func (s GrpcServer) Status(ctx context.Context, req *votingpb.StatusRequest) (*votingpb.StatusResponse, error) {
	status, err := s.service.Status(req.ShortId, req.ServerId)
	if err != nil {
		return nil, toStatusErr(err)
	}

	reply := &votingpb.StatusResponse{
		Poll:       toPbPoll(status.Poll),
		VoterCount: status.VoterCount,
	}

	return reply, nil

}

func (s GrpcServer) Voters(req *votingpb.VotersRequest, stream votingpb.VotingService_VotersServer) error {
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

func (s GrpcServer) Vote(ctx context.Context, req *votingpb.VoteRequest) (*votingpb.VoteResponse, error) {
	voteRequest := VoteRequest{
		ShortId:  req.ShortId,
		ServerId: req.ServerId,
		Voter:    toVoter(req.Voter),
		Options:  toBallotOptions(req.Options),
	}

	voteResponse, err := s.service.Vote(voteRequest)
	if err != nil {
		return nil, toStatusErr(err)
	}

	return toPbVoteResponse(voteResponse), nil
}

func (s GrpcServer) Count(ctx context.Context, req *votingpb.CountRequest) (*votingpb.CountResponse, error) {
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

	return toPbCountResponse(countResult), nil
}

func (s GrpcServer) Export(ctx context.Context, req *votingpb.ExportRequest) (*votingpb.ExportResponse, error) {
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
	return toPbExportResponse(result), nil
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

func toPbVoteResponse(voteResponse VoteReply) *votingpb.VoteResponse {
	pbVoteResponse := &votingpb.VoteResponse{
		Success: voteResponse.Success,
		Message: voteResponse.Message,
	}

	for _, voteResponseOption := range voteResponse.Options {
		pbVoteResponse.Options = append(pbVoteResponse.Options, toPbVoteResponseOption(voteResponseOption))
	}

	return pbVoteResponse
}

func toPbVoteResponseOption(voteResponseOption VoteReplyOption) *votingpb.VoteResponseOption {
	pbVoteResponseOption := &votingpb.VoteResponseOption{
		Rank:   voteResponseOption.Rank,
		Option: toPbOption(voteResponseOption.Option),
	}

	return pbVoteResponseOption
}

func toPbVoter(voter Voter) *votingpb.VotersResponse {
	pbVoter := &votingpb.VotersResponse{
		Id:         voter.Id,
		Username:   voter.Username,
		ExternalId: voter.ExternalId,
	}

	return pbVoter
}

func toPbCountResponse(countResult CountResult) *votingpb.CountResponse {
	countResponse := &votingpb.CountResponse{
		Success: countResult.Success,
		Message: countResult.Message,
		Poll:    toPbPoll(countResult.Poll),
		Method:  countResult.Method,
		Events:  toPbCountEvents(countResult.Events),
		Rounds:  toPbRoundSummaries(countResult.Rounds),
	}

	for _, elected := range countResult.Elected {
		countResponse.Elected = append(countResponse.Elected, toPbVoteResponseOption(elected))
	}

	return countResponse
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

func toPbExportResponse(exportResult ExportResult) *votingpb.ExportResponse {
	exportResponse := &votingpb.ExportResponse{
		Content: exportResult.Content,
		Poll:    toPbPoll(exportResult.Poll),
		Method:  exportResult.Method,
	}

	return exportResponse
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

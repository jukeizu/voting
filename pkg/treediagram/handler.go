package treediagram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/selection/api/protobuf-spec/selectionpb"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	"github.com/rs/zerolog"
)

const AppId = "intent.endpoint.voting"

type Handler struct {
	logger          zerolog.Logger
	client          votingpb.VotingClient
	selectionClient selectionpb.SelectionClient
	httpServer      *http.Server
}

func NewHandler(logger zerolog.Logger, client votingpb.VotingClient, selectionClient selectionpb.SelectionClient, addr string) Handler {
	logger = logger.With().Str("component", AppId).Logger()

	httpServer := http.Server{
		Addr: addr,
	}

	return Handler{logger, client, selectionClient, &httpServer}
}

func (h Handler) CreatePoll(request contract.Request) (*contract.Response, error) {
	createPollRequest, err := ParseCreatePollRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	reply, err := h.client.CreatePoll(context.Background(), createPollRequest)
	if err != nil {
		return FormatClientError(err)
	}

	message := &contract.Message{
		Embed: FormatNewPollReply(reply.Poll),
	}

	return &contract.Response{Messages: []*contract.Message{message}}, nil
}

func (h Handler) Poll(request contract.Request) (*contract.Response, error) {
	parsedPollRequest, err := ParsePollRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	pollReply, err := h.client.Poll(context.Background(), parsedPollRequest.pollRequest)
	if err != nil {
		return FormatClientError(err)
	}

	if pollReply.Poll.HasEnded {
		return contract.StringResponse("poll has ended"), nil
	}

	selectionRequest := &selectionpb.CreateSelectionRequest{
		AppId:      AppId + ".poll",
		InstanceId: pollReply.Poll.Id,
		UserId:     request.Author.Id,
		ServerId:   request.ServerId,
		Randomize:  true,
		BatchSize:  10,
		SortMethod: parsedPollRequest.sortMethod,
	}

	for _, option := range pollReply.Poll.Options {
		selectionOption := &selectionpb.Option{
			OptionId: option.Id,
			Content:  option.Content,
			Metadata: map[string]string{"url": option.Url},
		}

		selectionRequest.Options = append(selectionRequest.Options, selectionOption)
	}

	selection, err := h.selectionClient.CreateSelection(context.Background(), selectionRequest)
	if err != nil {
		return FormatClientError(err)
	}

	redirect := &contract.Message{
		Content:    fmt.Sprintf("<@!%s> The poll has been sent to your direct messages.", request.Author.Id),
		IsRedirect: true,
	}

	message := &contract.Message{
		Embed:            FormatPollReply(pollReply.Poll, selection),
		IsPrivateMessage: true,
	}

	return &contract.Response{Messages: []*contract.Message{redirect, message}}, nil
}

func (h Handler) PollStatus(request contract.Request) (*contract.Response, error) {
	req, err := ParsePollStatusRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	status, err := h.client.Status(context.Background(), req)
	if err != nil {
		return FormatClientError(err)
	}

	voters := []*votingpb.Voter{}
	if status.VoterCount <= 30 {
		v, err := h.voters(req.ShortId, req.ServerId)
		if err != nil {
			return FormatClientError(err)
		}

		voters = v
	}
	message := &contract.Message{
		Embed: FormatPollStatusReply(status, voters),
	}

	return &contract.Response{Messages: []*contract.Message{message}}, nil
}

func (h Handler) PollEnd(request contract.Request) (*contract.Response, error) {
	req, err := ParseEndPollRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	endPollReply, err := h.client.EndPoll(context.Background(), req)
	if err != nil {
		return FormatClientError(err)
	}

	return contract.StringResponse(FormatEndPollReply(endPollReply)), nil
}

func (h Handler) PollOpen(request contract.Request) (*contract.Response, error) {
	req, err := ParseOpenPollRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	openPollReply, err := h.client.OpenPoll(context.Background(), req)
	if err != nil {
		return FormatClientError(err)
	}

	message := &contract.Message{
		Embed: FormatOpenPollReply(openPollReply),
	}

	return &contract.Response{Messages: []*contract.Message{message}}, nil
}

func (h Handler) Vote(request contract.Request) (*contract.Response, error) {
	voterPollRequest := &votingpb.VoterPollRequest{
		VoterId:  request.Author.Id,
		ServerId: request.ServerId,
	}

	pollReply, err := h.client.VoterPoll(context.Background(), voterPollRequest)
	if err != nil {
		return FormatClientError(err)
	}

	parseSelectionRequest := &selectionpb.ParseSelectionRequest{
		AppId:      AppId + ".poll",
		InstanceId: pollReply.Poll.Id,
		UserId:     request.Author.Id,
		ServerId:   request.ServerId,
		Content:    request.Content,
	}

	parseSelectionReply, err := h.selectionClient.ParseSelection(context.Background(), parseSelectionRequest)
	if err != nil {
		return FormatClientError(err)
	}

	voteRequest := &votingpb.VoteRequest{
		ServerId: request.ServerId,
		Voter: &votingpb.Voter{
			Id:       request.Author.Id,
			Username: request.Author.Name,
		},
	}

	if request.Author.Discriminator != "" {
		voteRequest.Voter.Username += "#" + request.Author.Discriminator
	}

	for _, rankedOption := range parseSelectionReply.RankedOptions {
		ballotOption := &votingpb.BallotOption{
			Rank:     rankedOption.Rank,
			OptionId: rankedOption.Option.OptionId,
		}

		voteRequest.Options = append(voteRequest.Options, ballotOption)
	}

	voteReply, err := h.client.Vote(context.Background(), voteRequest)
	if err != nil {
		return FormatClientError(err)
	}

	message := &contract.Message{
		Content:          FormatVoteReply(pollReply.Poll, voteReply),
		IsPrivateMessage: true,
	}

	return &contract.Response{Messages: []*contract.Message{message}}, nil
}

func (h Handler) Count(request contract.Request) (*contract.Response, error) {
	countRequest, err := ParseCountRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	countReply, err := h.client.Count(context.Background(), countRequest)
	if err != nil {
		return FormatClientError(err)
	}

	message := FormatCountResult(countReply)

	return &contract.Response{Messages: []*contract.Message{message}}, nil
}

func (h Handler) voters(shortId string, serverId string) ([]*votingpb.Voter, error) {
	voters := []*votingpb.Voter{}

	votersRequest := &votingpb.VotersRequest{
		ShortId:  shortId,
		ServerId: serverId,
	}

	stream, err := h.client.Voters(context.Background(), votersRequest)
	if err != nil {
		return voters, err
	}

	for {
		voter, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return voters, err
		}

		voters = append(voters, voter)
	}

	return voters, nil
}

func (h Handler) Start() error {
	h.logger.Info().Msg("starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/createpoll", h.makeLoggingHttpHandlerFunc("createpoll", h.CreatePoll))
	mux.HandleFunc("/poll", h.makeLoggingHttpHandlerFunc("poll", h.Poll))
	mux.HandleFunc("/pollstatus", h.makeLoggingHttpHandlerFunc("pollstatus", h.PollStatus))
	mux.HandleFunc("/pollend", h.makeLoggingHttpHandlerFunc("pollend", h.PollEnd))
	mux.HandleFunc("/pollopen", h.makeLoggingHttpHandlerFunc("pollopen", h.PollOpen))
	mux.HandleFunc("/vote", h.makeLoggingHttpHandlerFunc("vote", h.Vote))
	mux.HandleFunc("/electioncount", h.makeLoggingHttpHandlerFunc("electioncount", h.Count))

	h.httpServer.Handler = mux

	return h.httpServer.ListenAndServe()
}

func (h Handler) Stop() error {
	h.logger.Info().Msg("stopping")

	return h.httpServer.Shutdown(context.Background())
}

func (h Handler) makeLoggingHttpHandlerFunc(name string, f func(contract.Request) (*contract.Response, error)) http.HandlerFunc {
	contractHandlerFunc := contract.MakeRequestHttpHandlerFunc(f)

	return func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			h.logger.Info().
				Str("intent", name).
				Str("took", time.Since(begin).String()).
				Msg("called")
		}(time.Now())

		contractHandlerFunc.ServeHTTP(w, r)
	}
}

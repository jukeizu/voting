package treediagram

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/selection/api/protobuf-spec/selectionpb"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	"github.com/jukeizu/voting/pkg/voting"
	"github.com/rs/zerolog"
)

const AppId = "intent.endpoint.voting"
const DefaultStatusMaxNumToElect = 4

type Handler struct {
	logger          zerolog.Logger
	client          votingpb.VotingServiceClient
	selectionClient selectionpb.SelectionClient
	httpServer      *http.Server
}

func NewHandler(logger zerolog.Logger, client votingpb.VotingServiceClient, selectionClient selectionpb.SelectionClient, addr string) Handler {
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

	message := FormatNewPollReply(reply.Poll)

	return &contract.Response{Messages: []*contract.Message{message}}, nil
}

func (h Handler) Poll(request contract.Request) (*contract.Response, error) {
	parsedPollRequest, err := ParsePollRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	return h.poll(parsedPollRequest.pollRequest)
}

func (h Handler) InteractionPoll(request contract.Interaction) (*contract.Response, error) {
	shortId := ParsePollShortId(request)

	req := &votingpb.PollRequest{
		ShortId:  shortId,
		ServerId: request.ServerId,
		VoterId:  request.User.Id,
	}

	return h.poll(req)
}

func (h Handler) PollStatus(request contract.Request) (*contract.Response, error) {
	req, err := ParsePollStatusRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	return h.pollStatus(req.ShortId, req.ServerId, false)
}

func (h Handler) InteractionPollStatus(request contract.Interaction) (*contract.Response, error) {
	shortId := ParsePollShortId(request)
	h.logger.Info().Msg(shortId)

	return h.pollStatus(shortId, request.ServerId, true)
}

func (h Handler) PollStatusRefresh(request contract.Interaction) (*contract.Response, error) {
	shortId := ParsePollShortId(request)

	return h.pollStatusEdit(shortId, request.ServerId, request.MessageId, false)
}

func (h Handler) PollEnd(request contract.Request) (*contract.Response, error) {
	req, err := ParseEndPollRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	_, err = h.client.EndPoll(context.Background(), req)
	if err != nil {
		return FormatClientError(err)
	}

	return h.pollStatus(req.ShortId, req.ServerId, false)
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

	message := FormatOpenPollReply(openPollReply)

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

	if pollReply.Poll.HasEnded {
		return FormatParseError(ParseError{Message: voting.ErrPollHasEnded.Message})
	}

	options, err := ParseVoteRequest(request, int(pollReply.Poll.AllowedUniqueVotes))
	if err != nil {
		return FormatParseError(err)
	}

	parseSelectionRequest := &selectionpb.ParseSelectionRequest{
		AppId:      AppId + ".poll",
		InstanceId: pollReply.Poll.Id,
		UserId:     request.Author.Id,
		ServerId:   request.ServerId,
		Content:    options,
	}

	parseSelectionReply, err := h.selectionClient.ParseSelection(context.Background(), parseSelectionRequest)
	if err != nil {
		return FormatClientErrorWithMessage(err, FormatVoteHelp(pollReply.Poll.AllowedUniqueVotes))
	}

	voteRequest := &votingpb.VoteRequest{
		ServerId: request.ServerId,
		Voter: &votingpb.Voter{
			ExternalId: request.Author.Id,
			Username:   request.Author.Name,
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
		Embed:            FormatVoteReply(pollReply.Poll, voteReply),
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

func (h Handler) Export(request contract.Request) (*contract.Response, error) {
	exportRequest, err := ParseExportRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	exportReply, err := h.client.Export(context.Background(), exportRequest)
	if err != nil {
		return FormatClientError(err)
	}

	message := FormatExportResult(exportReply)

	return &contract.Response{Messages: []*contract.Message{message}}, nil
}

func (h Handler) poll(req *votingpb.PollRequest) (*contract.Response, error) {
	pollReply, err := h.client.Poll(context.Background(), req)
	if err != nil {
		return FormatClientError(err)
	}

	if pollReply.Poll.HasEnded {
		return h.pollStatus(pollReply.Poll.ShortId, pollReply.Poll.ServerId, true)
	}

	selectionRequest := &selectionpb.CreateSelectionRequest{
		AppId:      AppId + ".poll",
		InstanceId: pollReply.Poll.Id,
		UserId:     req.VoterId,
		ServerId:   req.ServerId,
		Randomize:  true,
		BatchSize:  10,
		SortMethod: "number",
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

	message := &contract.Message{
		Embed:            FormatPollReply(pollReply.Poll, selection),
		IsPrivateMessage: true,
	}

	return &contract.Response{Messages: []*contract.Message{message}}, nil
}

func (h Handler) pollStatus(shortID string, serverID string, private bool) (*contract.Response, error) {
	return h.pollStatusEdit(shortID, serverID, "", private)
}

func (h Handler) pollStatusEdit(shortID string, serverID string, messageId string, private bool) (*contract.Response, error) {
	req := &votingpb.StatusRequest{
		ShortId:  shortID,
		ServerId: serverID,
	}

	status, err := h.client.Status(context.Background(), req)
	if err != nil {
		return FormatClientError(err)
	}

	voters := []*votingpb.VotersResponse{}
	if status.VoterCount <= 30 {
		v, err := h.voters(req.ShortId, req.ServerId)
		if err != nil {
			return FormatClientError(err)
		}

		voters = v
	}

	numToElect := status.Poll.AllowedUniqueVotes
	if numToElect > DefaultStatusMaxNumToElect {
		numToElect = DefaultStatusMaxNumToElect
	}

	countRequest := &votingpb.CountRequest{
		ShortId:    req.ShortId,
		ServerId:   req.ServerId,
		NumToElect: numToElect,
		Method:     "meekstv",
	}

	countReply, _ := h.client.Count(context.Background(), countRequest)

	message := FormatPollStatusReply(status, voters, countReply, messageId, private)

	return &contract.Response{Messages: []*contract.Message{message}}, nil
}

func (h Handler) voters(shortId string, serverId string) ([]*votingpb.VotersResponse, error) {
	voters := []*votingpb.VotersResponse{}

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
	mux.HandleFunc("/interactionpoll", h.makeLoggingInteractionHttpHandlerFunc("interactionpoll", h.InteractionPoll))
	mux.HandleFunc("/pollstatus", h.makeLoggingHttpHandlerFunc("pollstatus", h.PollStatus))
	mux.HandleFunc("/interactionpollstatus", h.makeLoggingInteractionHttpHandlerFunc("interactionpollstatus", h.InteractionPollStatus))
	mux.HandleFunc("/pollstatusrefresh", h.makeLoggingInteractionHttpHandlerFunc("pollstatusrefresh", h.PollStatusRefresh))
	mux.HandleFunc("/pollend", h.makeLoggingHttpHandlerFunc("pollend", h.PollEnd))
	mux.HandleFunc("/pollopen", h.makeLoggingHttpHandlerFunc("pollopen", h.PollOpen))
	mux.HandleFunc("/vote", h.makeLoggingHttpHandlerFunc("vote", h.Vote))
	mux.HandleFunc("/electioncount", h.makeLoggingHttpHandlerFunc("electioncount", h.Count))
	mux.HandleFunc("/electionexport", h.makeLoggingHttpHandlerFunc("electionexport", h.Export))

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

func (h Handler) makeLoggingInteractionHttpHandlerFunc(name string, f func(contract.Interaction) (*contract.Response, error)) http.HandlerFunc {
	contractHandlerFunc := contract.MakeInteractionHttpHandlerFunc(f)

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

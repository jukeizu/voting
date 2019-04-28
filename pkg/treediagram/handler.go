package treediagram

import (
	"context"
	"fmt"
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

	return contract.StringResponse(FormatNewPollReply(reply.Poll)), nil
}

func (h Handler) Poll(request contract.Request) (*contract.Response, error) {
	req, err := ParsePollRequest(request)
	if err != nil {
		return FormatParseError(err)
	}

	pollReply, err := h.client.Poll(context.Background(), req)
	if err != nil {
		return FormatClientError(err)
	}

	selectionRequest := &selectionpb.CreateSelectionRequest{
		AppId:      AppId + ".poll",
		InstanceId: pollReply.Poll.Id,
		UserId:     request.Author.Id,
		ServerId:   request.ServerId,
		Randomize:  true,
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
		Content:          FormatPollReply(pollReply.Poll, selection),
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

	return contract.StringResponse(FormatPollStatusReply(status)), nil
}

func (h Handler) Start() error {
	h.logger.Info().Msg("starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/createpoll", h.makeLoggingHttpHandlerFunc("createpoll", h.CreatePoll))
	mux.HandleFunc("/poll", h.makeLoggingHttpHandlerFunc("poll", h.Poll))
	mux.HandleFunc("/pollstatus", h.makeLoggingHttpHandlerFunc("pollstatus", h.PollStatus))

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

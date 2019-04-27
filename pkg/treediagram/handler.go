package treediagram

import (
	"context"
	"net/http"
	"time"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	"github.com/rs/zerolog"
)

type Handler struct {
	logger     zerolog.Logger
	client     votingpb.VotingClient
	httpServer *http.Server
}

func NewHandler(logger zerolog.Logger, client votingpb.VotingClient, addr string) Handler {
	logger = logger.With().Str("component", "intent.endpoint.voting").Logger()

	httpServer := http.Server{
		Addr: addr,
	}

	return Handler{logger, client, &httpServer}
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

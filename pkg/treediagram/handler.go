package treediagram

import (
	"context"
	"net/http"
	"time"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	createPollRequest, msg, err := ParseCreatePollRequest(request)
	if msg != "" {
		return contract.StringResponse(msg), nil
	}
	if err != nil {
		return nil, err
	}

	reply, err := h.client.CreatePoll(context.Background(), createPollRequest)
	if err != nil {
		return h.checkValidationError(err)
	}

	return contract.StringResponse(FormatNewPollReply(reply.Poll)), nil
}

func (h Handler) PollStatus(request contract.Request) (*contract.Response, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	req := &votingpb.StatusRequest{
		ServerId: request.ServerId,
	}

	if len(args) > 1 {
		req.ShortId = args[len(args)-1]
	}

	status, err := h.client.Status(context.Background(), req)
	if err != nil {
		return h.checkValidationError(err)
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

func (h Handler) checkValidationError(err error) (*contract.Response, error) {
	st, ok := status.FromError(err)
	if !ok {
		return nil, err
	}

	if st.Code() == codes.InvalidArgument {
		return contract.StringResponse(st.Message()), nil
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return contract.StringResponse(st.Message()), nil
	case codes.NotFound:
		return contract.StringResponse(st.Message()), nil
	}

	return nil, err
}

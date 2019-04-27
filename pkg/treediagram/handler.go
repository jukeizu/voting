package treediagram

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
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

	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf(":ballot_box: **A new poll has started** `%s`\n", reply.Poll.ShortId))

	if reply.Poll.Title != "" {
		buffer.WriteString(fmt.Sprintf("\n**%s**\n", reply.Poll.Title))
	}

	buffer.WriteString(fmt.Sprintf("\nType `!poll` to view the poll. A previous poll can be viewed via id. e.g. `!poll %s`", reply.Poll.ShortId))

	return contract.StringResponse(buffer.String()), nil
}

func (h Handler) Start() error {
	h.logger.Info().Msg("starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/createpoll", h.makeLoggingHttpHandlerFunc("createpoll", h.CreatePoll))

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

	return nil, err
}

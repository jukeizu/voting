package startup

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/jukeizu/voting/poll"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type PollStartup struct {
	logger     zerolog.Logger
	repository poll.Repository
}

func NewPollStartup(logger zerolog.Logger, dbAddress string) (PollStartup, error) {
	repository, err := poll.NewRepository(dbAddress)
	if err != nil {
		return PollStartup{}, err
	}

	return PollStartup{logger, repository}, nil
}

func (s PollStartup) Migrate() error {
	return s.repository.Migrate()
}

func (s PollStartup) RegisterServer(grpcServer *grpc.Server) {
	pollHandler := poll.NewPollHandler(s.logger, s.repository)
	optionsHandler := poll.NewOptionsHandler(s.logger, s.repository)

	pollServer := poll.NewServer(pollHandler, optionsHandler)

	pollpb.RegisterPollsServer(grpcServer, pollServer)
}

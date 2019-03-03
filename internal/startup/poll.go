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
	createPollCommandHandler := poll.NewCreatePollCommandHandler(s.logger, s.repository)
	pollQueryHandler := poll.NewPollQueryHandler(s.logger, s.repository)
	optionsQueryHandler := poll.NewOptionsQueryHandler(s.logger, s.repository)
	endPollCommandHandler := poll.NewEndPollCommandHandler(s.logger, s.repository)

	pollServer := poll.NewServer(
		createPollCommandHandler,
		pollQueryHandler,
		optionsQueryHandler,
		endPollCommandHandler,
	)

	pollpb.RegisterPollsServer(grpcServer, pollServer)
}

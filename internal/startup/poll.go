package startup

import (
	"github.com/jukeizu/voting/poll"
	"github.com/rs/zerolog"
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

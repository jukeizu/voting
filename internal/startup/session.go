package startup

import (
	"github.com/jukeizu/voting/pkg/session"
	"github.com/rs/zerolog"
)

type SessionStartup struct {
	logger     zerolog.Logger
	repository session.Repository
}

func NewSessionStartup(logger zerolog.Logger, dbAddress string) (SessionStartup, error) {
	repository, err := session.NewRepository(dbAddress)
	if err != nil {
		return SessionStartup{}, err
	}

	return SessionStartup{logger, repository}, nil
}

func (s SessionStartup) Migrate() error {
	return s.repository.Migrate()
}

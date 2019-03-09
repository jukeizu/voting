package startup

import (
	"github.com/jukeizu/voting/api/protobuf-spec/sessionpb"
	"github.com/jukeizu/voting/session"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
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

func (s SessionStartup) RegisterServer(grpcServer *grpc.Server) {
	service := session.NewDefaultService(s.logger, s.repository)
	server := session.NewGrpcServer(service)

	sessionpb.RegisterSessionServer(grpcServer, server)
}

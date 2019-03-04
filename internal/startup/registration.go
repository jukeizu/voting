package startup

import (
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/jukeizu/voting/registration"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type RegistrationStartup struct {
	logger     zerolog.Logger
	repository registration.Repository
}

func NewRegistrationStartup(logger zerolog.Logger, dbAddress string) (RegistrationStartup, error) {
	repository, err := registration.NewRepository(dbAddress)
	if err != nil {
		return RegistrationStartup{}, err
	}

	return RegistrationStartup{logger, repository}, nil
}

func (s RegistrationStartup) Migrate() error {
	return s.repository.Migrate()
}

func (s RegistrationStartup) RegisterServer(grpcServer *grpc.Server) {
	registrationHandler := registration.NewRegistrationHandler(s.logger, s.repository)

	registrationServer := registration.NewServer(registrationHandler)

	registrationpb.RegisterRegistrationServer(grpcServer, registrationServer)
}

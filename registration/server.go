package registration

import (
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/rs/zerolog"
)

type Server struct {
	logger     zerolog.Logger
	repository Repository
}

func NewServer(logger zerolog.Logger, repository Repository) Server {
	return Server{logger, repository}
}

func (s Server) RegisterUser(req *registrationpb.RegisterUserRequest) (*registrationpb.RegisterUserReply, error) {
	user, err := s.repository.SaveUser(req.ExternalId, req.Username)
	if err != nil {
		return nil, err
	}

	s.logger.Info().
		Str("externalId", user.ExternalId).
		Str("username", user.Username).
		Msg("registered user")

	return &registrationpb.RegisterUserReply{User: user}, nil
}

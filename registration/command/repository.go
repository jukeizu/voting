package registration

import "github.com/jukeizu/voting/api/protobuf-spec/registrationpb"

type Repository interface {
	RegisterVoter(externalId, username string, canVote bool) (*registrationpb.Voter, error)
}

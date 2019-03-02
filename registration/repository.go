package registration

import "github.com/jukeizu/voting/domain/entities"

type Repository interface {
	RegisterVoter(externalId, username string, canVote bool) (*entities.Voter, error)
}

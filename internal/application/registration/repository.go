package registration

import (
	"github.com/jukeizu/voting/internal/application"
)

type Repository struct {
	db application.Database
}

func NewRepository(db application.Database) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) SaveOrganization(org application.Organization) (application.Organization, error) {
	q := `INSERT INTO organization (name, external_id)
		VALUES ($1, $2)
		ON CONFLICT (name, external_id)
		DO UPDATE SET updated = now()
		RETURNING id, max_concurrent_polls`

	err := r.db.
		QueryRow(q, org.Name, org.ExternalId).
		Scan(&org.Id, &org.MaxConcurrentPolls)

	return org, err
}

func (r Repository) SaveVoter(voter application.Voter) (application.Voter, error) {
	q := `INSERT INTO voter (organization_id, external_id, name)
		VALUES ($1, $2, $3)
		ON CONFLICT (organization_id, external_id)
		DO UPDATE SET name = excluded.name, updated = now()
		RETURNING id, can_vote`

	err := r.db.
		QueryRow(q, voter.Organization.Id, voter.ExternalId, voter.Name).
		Scan(&voter.ID, &voter.CanVote)

	return voter, err
}

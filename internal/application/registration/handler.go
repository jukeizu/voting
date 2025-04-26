package registration

import (
	"github.com/jukeizu/voting/internal/application"
)

type Handler interface {
	RegisterIdentity(application.Identity) (application.Voter, error)
}

var _ Handler = &handler{}

type handler struct {
	r Repository
}

func NewHandler(r Repository) Handler {
	return &handler{
		r: r,
	}
}

func (h handler) RegisterIdentity(identity application.Identity) (application.Voter, error) {
	org, err := h.r.SaveOrganization(application.Organization{
		Name:       identity.OrganizationName,
		ExternalId: identity.OrganizationExternalId,
	})
	if err != nil {
		return application.Voter{}, err
	}

	return h.r.SaveVoter(application.Voter{
		ExternalId:   identity.VoterExternalId,
		Name:         identity.VoterName,
		Organization: org,
	})
}

package registration

import (
	"github.com/jukeizu/voting/internal/application"
)

type validator struct {
	h Handler
}

var _ Handler = &validator{}

func NewValidatingHandler(h Handler) Handler {
	return &validator{
		h: h,
	}
}

func (v *validator) RegisterIdentity(identity application.Identity) (application.Voter, error) {
	if len(identity.OrganizationName) < 1 {
		return application.Voter{}, application.ValidationError{
			Message: "Invalid identity: OrganizationName is empty.",
		}
	}

	if len(identity.OrganizationExternalId) < 1 {
		return application.Voter{}, application.ValidationError{
			Message: "Invalid identity: OrganizationExternalId is empty.",
		}
	}

	if len(identity.VoterExternalId) < 1 {
		return application.Voter{}, application.ValidationError{
			Message: "Invalid identity: VoterExternalId is empty.",
		}
	}

	return v.h.RegisterIdentity(identity)
}

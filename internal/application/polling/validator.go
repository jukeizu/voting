package polling

import (
	"time"

	"github.com/jukeizu/voting/internal/application"
)

type validator struct {
	h Handler
	r Repository
}

var _ Handler = &validator{}

func NewValidatingHandler(h Handler, r Repository) Handler {
	return &validator{
		h: h,
		r: r,
	}
}

func (v *validator) CreatePoll(req CreatePollRequest) (application.Poll, error) {
	canCreate, err := v.r.CanCreatePoll(req.Creator.Organization.Id)
	if err != nil {
		return application.Poll{}, err
	}

	if !canCreate {
		return application.Poll{}, application.ErrMaxConcurrent
	}

	if len(req.Candidates) < 1 {
		return application.Poll{}, application.ErrNoCandidates
	}

	if req.Expiration != nil && req.Expiration.UTC().Before(time.Now().UTC()) {
		return application.Poll{}, application.ErrPastPollExpiration
	}

	return v.h.CreatePoll(req)
}

func (v *validator) Poll(req PollRequest) (application.Poll, error) {
	return v.h.Poll(req)
}

func (v *validator) Open(req OpenPollRequest) (application.Poll, error) {
	err := v.validateOwner(req.PollId, req.Voter)
	if err != nil {
		return application.Poll{}, err
	}

	return v.h.Open(req)
}

func (v *validator) End(req EndPollRequest) (application.Poll, error) {
	err := v.validateOwner(req.PollId, req.Voter)
	if err != nil {
		return application.Poll{}, err
	}

	return v.h.End(req)
}

func (v *validator) PollIdOrLatest(pollID string, voter application.Voter) (string, error) {
	return v.h.PollIdOrLatest(pollID, voter)
}

func (v *validator) validateOwner(pollID string, voter application.Voter) error {
	pollID, err := v.h.PollIdOrLatest(pollID, voter)
	if err != nil {
		return err
	}

	creatorID, err := v.r.PollCreatorID(pollID)
	if err != nil {
		return err
	}

	if creatorID != voter.ID {
		return application.ErrNotOwner
	}

	return nil
}

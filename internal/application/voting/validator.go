package voting

import (
	"database/sql"

	"github.com/jukeizu/voting/internal/application"
	"github.com/jukeizu/voting/internal/application/polling"
)

type validator struct {
	h  Handler
	r  Repository
	pr polling.Repository
}

var _ Handler = &validator{}

func NewValidatingHandler(h Handler, r Repository, pr polling.Repository) Handler {
	return &validator{
		h:  h,
		r:  r,
		pr: pr,
	}
}

func (v *validator) Choices(req ChoicesRequest) (ChoicesResponse, error) {
	if !req.Voter.CanVote {
		return ChoicesResponse{}, application.ErrCannotVote
	}

	ended, err := v.pr.PollHasEnded(req.PollID)
	if err != nil {
		return ChoicesResponse{}, err
	}

	if ended {
		return ChoicesResponse{}, application.ErrPollHasEnded
	}

	return v.h.Choices(req)
}

func (v *validator) SubmitVote(req SubmitVoteRequest) (SubmitVoteResponse, error) {
	if !req.Voter.CanVote {
		return SubmitVoteResponse{}, application.ErrCannotVote
	}

	session, err := v.r.VoterSession(req.Voter.ID)
	if err != nil && err != sql.ErrNoRows {
		return SubmitVoteResponse{}, err
	}

	ended, err := v.pr.PollHasEnded(session.PollID)
	if err != nil {
		return SubmitVoteResponse{}, err
	}

	if ended {
		return SubmitVoteResponse{}, application.ErrPollHasEnded
	}

	choices, err := v.r.BallotChoices(session)
	if err != nil && err != sql.ErrNoRows {
		return SubmitVoteResponse{}, err
	}

	err = v.validateBallotOptions(req, len(choices))
	if err != nil {
		return SubmitVoteResponse{}, err
	}

	return v.h.SubmitVote(req)
}

func (v *validator) Ballot(req BallotRequest) (BallotResponse, error) {
	return v.h.Ballot(req)
}

func (v *validator) validateBallotOptions(req SubmitVoteRequest, max int) error {
	if len(req.Options) < 1 {
		return application.ErrNoCandidates
	}

	dupeCheck := map[int32]bool{}

	for _, ballotOption := range req.Options {
		// Verify in expected range.
		if ballotOption.Number > int32(max) || ballotOption.Number < 1 {
			return application.ErrInvalidOrDuplicateOptions
		}
		// Check for duplicate.
		_, found := dupeCheck[ballotOption.Number]
		if found {
			return application.ErrInvalidOrDuplicateOptions
		}
		dupeCheck[ballotOption.Number] = true
	}

	return nil
}

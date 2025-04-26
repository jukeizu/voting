package voting

import (
	"database/sql"

	"github.com/jukeizu/voting/internal/application"
	"github.com/jukeizu/voting/internal/application/polling"
)

type ChoicesRequest struct {
	Voter  application.Voter
	PollID string
}

type ChoicesResponse struct {
	Title   string
	Choices []application.Choice
}

type SubmitVoteRequest struct {
	PollID  string
	Voter   application.Voter
	Options []application.BallotOption
}

type SubmitVoteResponse struct {
	RankedChoices []application.RankedChoice
}

type BallotRequest struct {
	PollID string
	Voter  application.Voter
}

type BallotResponse struct {
	RankedChoices []application.RankedChoice
}

type Handler interface {
	Choices(ChoicesRequest) (ChoicesResponse, error)
	SubmitVote(SubmitVoteRequest) (SubmitVoteResponse, error)
	Ballot(req BallotRequest) (BallotResponse, error)
}

func NewHandler(r Repository, pr polling.Repository) Handler {
	return &handler{
		r:  r,
		pr: pr,
	}
}

type handler struct {
	r  Repository
	pr polling.Repository
}

func (h *handler) SubmitVote(req SubmitVoteRequest) (SubmitVoteResponse, error) {
	session, err := h.r.VoterSession(req.Voter.ID)
	if err != nil && err != sql.ErrNoRows {
		return SubmitVoteResponse{}, err
	}

	err = h.r.VoidPreviousBallot(req.Voter.ID, session.PollID)
	if err != nil && err != sql.ErrNoRows {
		return SubmitVoteResponse{}, err
	}

	availableChoices, err := h.r.BallotChoices(session)
	if err != nil {
		return SubmitVoteResponse{}, err
	}

	rankedChoices := make([]application.RankedChoice, len(req.Options))

	for i, ballotOption := range req.Options {
		rankedChoices[i] = application.RankedChoice{
			Rank:   ballotOption.Rank,
			Choice: availableChoices[ballotOption.Number-1],
		}
	}

	err = h.r.SaveRankedChoices(req.Voter.ID, rankedChoices)
	if err != nil {
		return SubmitVoteResponse{}, err
	}

	return SubmitVoteResponse{
		RankedChoices: rankedChoices,
	}, nil
}

func (h *handler) Choices(req ChoicesRequest) (ChoicesResponse, error) {
	session, err := h.r.SaveVoterSession(req.Voter.ID, req.PollID)
	if err != nil {
		return ChoicesResponse{}, err
	}

	choices, err := h.r.BallotChoices(session)
	if err != nil {
		return ChoicesResponse{}, err
	}

	title, err := h.pr.PollTitle(req.PollID)
	if err != nil {
		return ChoicesResponse{}, err
	}

	return ChoicesResponse{
		Title:   title,
		Choices: choices,
	}, nil
}

func (h *handler) Ballot(req BallotRequest) (BallotResponse, error) {
	session, err := h.r.VoterSession(req.Voter.ID)
	if err != nil {
		return BallotResponse{}, err
	}

	ballot, err := h.r.Ballot(session)
	if err != nil {
		return BallotResponse{}, err
	}

	return BallotResponse{
		RankedChoices: ballot.RankedChoices,
	}, nil
}

var _ Handler = &handler{}

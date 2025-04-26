package polling

import (
	"database/sql"
	"time"

	"github.com/jukeizu/voting/internal/application"
)

type PollRequest struct {
	ID    string
	Voter application.Voter
}

type CreatePollRequest struct {
	Creator       application.Voter
	Title         string
	Ranked        bool
	Expiration    *time.Time
	ManuallyEnded bool
	Candidates    []application.Candidate
}

type OpenPollRequest struct {
	PollId  string
	Expires *time.Time
	Voter   application.Voter
}

type EndPollRequest struct {
	PollId string
	Voter  application.Voter
}

type Handler interface {
	CreatePoll(CreatePollRequest) (application.Poll, error)
	Poll(PollRequest) (application.Poll, error)
	Open(OpenPollRequest) (application.Poll, error)
	End(EndPollRequest) (application.Poll, error)
	PollIdOrLatest(pollID string, voter application.Voter) (string, error)
}

func NewHandler(r Repository) Handler {
	return &handler{
		r: r,
	}
}

type handler struct {
	r Repository
}

var _ Handler = &handler{}

func (h *handler) CreatePoll(req CreatePollRequest) (application.Poll, error) {
	return h.r.SavePoll(req)
}

func (h *handler) Poll(req PollRequest) (application.Poll, error) {
	pollID, err := h.PollIdOrLatest(req.ID, req.Voter)
	if err != nil {
		return application.Poll{}, err
	}

	poll, err := h.r.Poll(pollID)
	if err == sql.ErrNoRows {
		return application.Poll{}, application.NotFound{}
	}

	return poll, err
}

func (h *handler) Open(req OpenPollRequest) (application.Poll, error) {
	pollID, err := h.PollIdOrLatest(req.PollId, req.Voter)
	if err != nil {
		return application.Poll{}, err
	}

	err = h.r.Open(pollID, req.Expires)
	if err != nil {
		return application.Poll{}, err
	}

	return h.Poll(PollRequest{
		ID: pollID,
	})
}

func (h *handler) End(req EndPollRequest) (application.Poll, error) {
	pollID, err := h.PollIdOrLatest(req.PollId, req.Voter)
	if err != nil {
		return application.Poll{}, err
	}

	err = h.r.End(pollID)
	if err != nil {
		return application.Poll{}, err
	}

	return h.Poll(PollRequest{
		ID: pollID,
	})
}

func (h *handler) PollIdOrLatest(pollID string, voter application.Voter) (string, error) {
	if len(pollID) > 0 {
		return pollID, nil
	}

	id, err := h.r.LatestPollID(voter.Organization.Id)
	if err == sql.ErrNoRows {
		return "", application.NotFound{}
	}

	return id, err
}

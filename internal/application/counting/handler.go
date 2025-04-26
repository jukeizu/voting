package counting

import (
	"bytes"
	"fmt"
	"runtime/debug"

	"github.com/jukeizu/voting/internal/application"
	"github.com/jukeizu/voting/internal/application/polling"
	"github.com/shawntoffel/meekstv"
)

type CountRequest struct {
	PollID string
	Method string
	Seats  int
}

type CountResponse struct {
	Title      string
	Candidates int
	Ballots    int
	Elected    []application.ElectedCandidate
	Report     string
	Method     string
}

type ExportRequest struct {
	PollID    string
	Seats     int
	Withdrawn []int
}

type ExportResponse struct {
	BLT string
}

type Handler interface {
	Count(CountRequest) (CountResponse, error)
	Export(ExportRequest) (ExportResponse, error)
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

var _ Handler = &handler{}

func (h *handler) Count(req CountRequest) (countResponse CountResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic while counting votes: %s", debug.Stack())
		}
	}()

	config, err := h.buildMeekSTVConfig(req.PollID, req.Seats)
	if err != nil {
		return CountResponse{}, err
	}

	result, err := meekstv.Count(config)
	if err != nil {
		return CountResponse{}, err
	}

	report := bytes.Buffer{}
	err = result.Detail.WriteReport(&report)
	if err != nil {
		return CountResponse{}, err
	}

	elected := make([]application.ElectedCandidate, len(result.Elected))
	for i, name := range result.Elected {
		elected[i] = application.ElectedCandidate{
			Rank: int32(i + 1),
			Name: name,
		}
	}

	title, err := h.pr.PollTitle(req.PollID)
	if err != nil {
		return CountResponse{}, err
	}

	countResponse = CountResponse{
		Title:      title,
		Candidates: len(config.Candidates),
		Ballots:    config.Ballots.TotalCount(),
		Method:     "MeekSTV",
		Elected:    elected,
		Report:     report.String(),
	}

	return
}

func (h *handler) Export(req ExportRequest) (ExportResponse, error) {
	title, err := h.pr.PollTitle(req.PollID)
	if err != nil {
		return ExportResponse{}, err
	}

	config, err := h.buildMeekSTVConfig(req.PollID, req.Seats)
	if err != nil {
		return ExportResponse{}, err
	}

	config.WithdrawnCandidates = req.Withdrawn

	return ExportResponse{
		BLT: ExportBlt(title, config),
	}, nil
}

func (h *handler) buildMeekSTVConfig(pollID string, seats int) (meekstv.Config, error) {
	candidates, err := h.r.Candidates(pollID)
	if err != nil {
		return meekstv.Config{}, err
	}

	ballots, err := h.ballots(pollID)
	if err != nil {
		return meekstv.Config{}, err
	}

	config := meekstv.Config{
		Seats:      seats,
		Candidates: candidates.Names(),
		Ballots:    make([]meekstv.Ballot, len(ballots)),
	}

	for i, ballot := range ballots {
		config.Ballots[i] = meekstv.Ballot{
			Count:       1,
			Preferences: ballot.FlatPreferences(),
		}
	}

	return config, nil
}

func (h *handler) ballots(pollID string) ([]application.Ballot, error) {
	ballots := []application.Ballot{}

	voterIds, err := h.r.VoterIDs(pollID)
	if err != nil {
		return []application.Ballot{}, err
	}

	for _, voterId := range voterIds {
		ballot, err := h.r.Ballot(pollID, voterId)
		if err != nil {
			return []application.Ballot{}, err
		}

		ballots = append(ballots, ballot)
	}

	return ballots, nil
}

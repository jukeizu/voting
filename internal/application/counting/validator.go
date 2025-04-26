package counting

import (
	"strings"

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

func (v *validator) Count(req CountRequest) (CountResponse, error) {
	ended, err := v.pr.PollHasEnded(req.PollID)
	if err != nil {
		return CountResponse{}, err
	}

	if !ended {
		return CountResponse{}, application.ErrPollHasNotEnded
	}

	if !strings.EqualFold(req.Method, "meekstv") {
		return CountResponse{}, application.ErrUnknownCountMethod(req.Method)
	}

	return v.h.Count(req)
}

func (v *validator) Export(req ExportRequest) (ExportResponse, error) {
	ended, err := v.pr.PollHasEnded(req.PollID)
	if err != nil {
		return ExportResponse{}, err
	}

	if !ended {
		return ExportResponse{}, application.ErrPollHasNotEnded
	}

	return v.h.Export(req)
}

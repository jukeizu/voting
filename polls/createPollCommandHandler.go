package polls

import "github.com/jukeizu/voting/domain/entities"

type CreatePollRequest struct {
	CreatorId            string
	Title                string
	AllowedUniqueOptions int32
	Options              []entities.Option
}

type CreatePollCommandHandler struct {
}

func (h CreatePollCommandHandler) Handle(req CreatePollRequest) (*entities.Poll, error) {
	return nil, nil
}

package polls

import "github.com/jukeizu/voting/domain/entities"

type Repository interface {
	CreatePoll(poll entities.Poll) (*entities.Poll, error)
	Poll(id string) (*entities.Poll, error)
	PollCreator(id string) (string, error)
	Options(pollId string) ([]entities.Option, error)
	EndPoll(id string) (*entities.Poll, error)
}

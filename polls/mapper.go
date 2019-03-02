package polls

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/jukeizu/voting/domain/entities"
)

func mapToPbPoll(poll *entities.Poll) *pollpb.Poll {
	pbPoll := &pollpb.Poll{
		Id:                 poll.Id,
		Title:              poll.Title,
		CreatorId:          poll.CreatorId,
		AllowedUniqueVotes: poll.AllowedUniqueVotes,
		HasEnded:           poll.HasEnded,
	}

	pbPoll.Options = mapToPbOptions(poll.Options)

	return pbPoll
}

func mapToPbOptions(options []entities.Option) []*pollpb.Option {
	pbOptions := []*pollpb.Option{}

	for _, option := range options {
		pbOption := &pollpb.Option{
			Id:      option.Id,
			PollId:  option.PollId,
			Content: option.Content,
		}

		pbOptions = append(pbOptions, pbOption)
	}

	return pbOptions
}

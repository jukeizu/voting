package polls

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/jukeizu/voting/domain/entities"
)

func mapToPb(poll *entities.Poll) *pollpb.Poll {
	pbPoll := &pollpb.Poll{
		Id:                 poll.Id,
		Title:              poll.Title,
		CreatorId:          poll.CreatorId,
		AllowedUniqueVotes: poll.AllowedUniqueVotes,
		HasEnded:           poll.HasEnded,
	}

	for _, option := range poll.Options {
		pbOption := &pollpb.Option{
			Id:      option.Id,
			PollId:  option.PollId,
			Content: option.Content,
		}

		pbPoll.Options = append(pbPoll.Options, pbOption)
	}

	return pbPoll
}

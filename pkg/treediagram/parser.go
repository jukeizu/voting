package treediagram

import (
	"bytes"
	"flag"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	shellwords "github.com/mattn/go-shellwords"
)

func ParseCreatePollRequest(request contract.Request) (*votingpb.CreatePollRequest, string, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, "", err
	}

	outputBuffer := bytes.NewBuffer([]byte{})

	parser := flag.NewFlagSet("pollnew", flag.ContinueOnError)
	parser.SetOutput(outputBuffer)

	title := parser.String("t", "", "The poll title")
	allowedUniqueVotes := parser.Int("n", 1, "The number of unique votes a user can submit.")

	err = parser.Parse(args[1:])
	if err != nil {
		return nil, outputBuffer.String(), err
	}

	createPollRequest := &votingpb.CreatePollRequest{
		Title:              *title,
		AllowedUniqueVotes: int32(*allowedUniqueVotes),
		ServerId:           request.ServerId,
		CreatorId:          request.Author.Id,
	}

	for _, content := range parser.Args() {
		option := &votingpb.Option{
			Content: content,
		}

		createPollRequest.Options = append(createPollRequest.Options, option)
	}

	return createPollRequest, "", nil
}

package treediagram

import (
	"bytes"
	"flag"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	shellwords "github.com/mattn/go-shellwords"
)

func ParseCreatePollRequest(request contract.Request) (*votingpb.CreatePollRequest, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	outputBuffer := bytes.NewBuffer([]byte{})

	parser := flag.NewFlagSet("pollnew", flag.ContinueOnError)
	parser.SetOutput(outputBuffer)

	title := parser.String("t", "", "The poll title")
	allowedUniqueVotes := parser.Int("n", 1, "The number of unique votes a user can submit.")

	err = parser.Parse(args[1:])
	if err != nil {
		return nil, ParseError{Message: outputBuffer.String()}
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

	return createPollRequest, nil
}

func ParsePollStatusRequest(request contract.Request) (*votingpb.StatusRequest, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	req := &votingpb.StatusRequest{
		ServerId: request.ServerId,
	}

	if len(args) > 1 {
		req.ShortId = args[len(args)-1]
	}

	return req, nil
}

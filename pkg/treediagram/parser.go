package treediagram

import (
	"bytes"
	"flag"
	"strings"

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

func ParsePollRequest(request contract.Request) (*votingpb.PollRequest, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	req := &votingpb.PollRequest{
		ServerId: request.ServerId,
	}

	if len(args) > 1 {
		req.ShortId = args[len(args)-1]
	}

	return req, nil
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

func ParseEndPollRequest(request contract.Request) (*votingpb.EndPollRequest, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	req := &votingpb.EndPollRequest{
		ServerId:    request.ServerId,
		RequesterId: request.Author.Id,
	}

	if len(args) > 1 {
		req.ShortId = args[len(args)-1]
	}

	return req, nil
}

func ParseCountRequest(request contract.Request) (*votingpb.CountRequest, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	outputBuffer := bytes.NewBuffer([]byte{})

	parser := flag.NewFlagSet("electioncount", flag.ContinueOnError)
	parser.SetOutput(outputBuffer)

	shortID := parser.String("id", "", "The poll id. Current poll if not specified.")
	method := parser.String("m", "meekstv", "The counting method.")
	numToElect := parser.Int("n", 1, "Number of seats to fill.")
	exclude := parser.String("exclude", "", "Comma delimited list of candidates to exclude before counting.")

	err = parser.Parse(args[1:])
	if err != nil {
		return nil, ParseError{Message: outputBuffer.String()}
	}

	coundRequest := &votingpb.CountRequest{
		ShortId:    *shortID,
		ServerId:   request.ServerId,
		NumToElect: int32(*numToElect),
		Method:     *method,
		ToExclude:  strings.Split(*exclude, ","),
	}

	return coundRequest, nil
}

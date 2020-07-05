package treediagram

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	shellwords "github.com/mattn/go-shellwords"
)

type ParsedPollRequest struct {
	pollRequest *votingpb.PollRequest
	sortMethod  string
}

func ParseCreatePollRequest(request contract.Request) (*votingpb.CreatePollRequest, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	outputBuffer := bytes.NewBuffer([]byte{})

	parser := flag.NewFlagSet("pollnew", flag.ContinueOnError)
	parser.SetOutput(outputBuffer)

	endsFormatExample := time.Now().UTC().Add(time.Hour * 36).Format("1/2/06 15:04")

	title := parser.String("t", "", "The poll title")
	allowedUniqueVotes := parser.Int("n", 1, "The number of unique votes a user can submit.")
	ends := parser.String("ends", "", fmt.Sprintf("The UTC end time for the poll such as '%s' (format \"M/d/yy H:mm\")", endsFormatExample))

	err = parser.Parse(args[1:])
	if err != nil {
		return nil, ParseError{Message: outputBuffer.String()}
	}

	t, err := parseEndTime("1/2/06 15:04", *ends)
	if err != nil {
		return nil, ParseError{Message: err.Error()}
	}

	if *ends != "" && t.Before(time.Now().UTC()) {
		return nil, ParseError{Message: "Poll end time must be in the future."}
	}

	createPollRequest := &votingpb.CreatePollRequest{
		Title:              *title,
		AllowedUniqueVotes: int32(*allowedUniqueVotes),
		ServerId:           request.ServerId,
		CreatorId:          request.Author.Id,
		Expires:            t.Unix(),
	}

	for _, content := range parser.Args() {
		option := &votingpb.Option{
			Content: content,
		}

		createPollRequest.Options = append(createPollRequest.Options, option)
	}

	return createPollRequest, nil
}

func ParsePollRequest(request contract.Request) (*ParsedPollRequest, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	outputBuffer := bytes.NewBuffer([]byte{})

	parser := flag.NewFlagSet("poll", flag.ContinueOnError)
	parser.SetOutput(outputBuffer)

	shortID := parser.String("id", "", "The poll id. Defaults to the most recent poll in a server.")
	sortMethod := parser.String("sort", "number", "Sort the poll by [abc, number]")

	err = parser.Parse(args[1:])
	if err != nil {
		return nil, ParseError{Message: outputBuffer.String()}
	}

	req := &votingpb.PollRequest{
		ServerId: request.ServerId,
		ShortId:  *shortID,
		VoterId:  request.Author.Id,
	}

	parsedSortMethod, err := parseSortMethod(*sortMethod)
	if err != nil {
		return nil, err
	}

	parsedPollRequest := &ParsedPollRequest{
		pollRequest: req,
		sortMethod:  parsedSortMethod,
	}

	return parsedPollRequest, nil
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

func ParseOpenPollRequest(request contract.Request) (*votingpb.OpenPollRequest, error) {
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	outputBuffer := bytes.NewBuffer([]byte{})

	parser := flag.NewFlagSet("pollopen", flag.ContinueOnError)
	parser.SetOutput(outputBuffer)

	endsFormatExample := time.Now().UTC().Add(time.Hour * 36).Format("1/2/06 15:04")

	shortID := parser.String("id", "", "The poll id. Defaults to the most recent poll in a server.")
	ends := parser.String("ends", "", fmt.Sprintf("The UTC end time for the poll such as '%s' (format \"M/d/yy H:mm\")", endsFormatExample))

	err = parser.Parse(args[1:])
	if err != nil {
		return nil, ParseError{Message: outputBuffer.String()}
	}

	t, err := parseEndTime("1/2/06 15:04", *ends)
	if err != nil {
		return nil, ParseError{Message: err.Error()}
	}

	if *ends != "" && t.Before(time.Now().UTC()) {
		return nil, ParseError{Message: "Poll end time must be in the future."}
	}

	req := &votingpb.OpenPollRequest{
		ShortId:     *shortID,
		ServerId:    request.ServerId,
		RequesterId: request.Author.Id,
		Expires:     t.Unix(),
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

	shortID := parser.String("id", "", "The poll id. Defaults to the most recent poll in a server.")
	method := parser.String("m", "meekstv", "The counting method.")
	numToElect := parser.Int("n", 1, "Number of seats to fill.")
	exclude := parser.String("exclude", "", "Comma delimited list of candidates to exclude before counting.")

	err = parser.Parse(args[1:])
	if err != nil {
		return nil, ParseError{Message: outputBuffer.String()}
	}

	countRequest := &votingpb.CountRequest{
		ShortId:    *shortID,
		ServerId:   request.ServerId,
		NumToElect: int32(*numToElect),
		Method:     *method,
		ToExclude:  strings.Split(*exclude, ","),
	}

	return countRequest, nil
}

func parseSortMethod(input string) (string, error) {
	sortMethodMap := map[string]string{
		"abc":    "alphabetical",
		"number": "number",
	}

	sortMethod, ok := sortMethodMap[strings.ToLower(input)]
	if !ok {
		return "", ParseError{Message: "invalid sort value"}
	}

	return sortMethod, nil
}

func parseEndTime(layout string, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	return time.Parse(layout, value)
}

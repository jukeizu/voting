package treediagram

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/selection/api/protobuf-spec/selectionpb"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FormatNewPollReply(poll *votingpb.Poll) string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf(":ballot_box: **A new poll has started** `%s`\n", poll.ShortId))

	if poll.Title != "" {
		buffer.WriteString(fmt.Sprintf("\n**%s**\n", poll.Title))
	}

	if poll.Expires > (time.Time{}).Unix() {
		formatedTime := time.Unix(poll.Expires, 0).UTC().Format("Jan 2, 2006 15:04:05 MST")

		buffer.WriteString(fmt.Sprintf("\nPoll ends `%s`\n", formatedTime))
	}

	buffer.WriteString(fmt.Sprintf("\nView the poll with `!poll` or `!poll -id %s`", poll.ShortId))

	return buffer.String()
}

func FormatPollStatusReply(status *votingpb.StatusReply, voters []*votingpb.Voter) string {
	buffer := bytes.Buffer{}

	buffer.WriteString(":ballot_box: ")
	if status.Poll.Title != "" {
		buffer.WriteString(fmt.Sprintf("**%s** ", status.Poll.Title))
	}

	buffer.WriteString(fmt.Sprintf("`%s`\n", status.Poll.ShortId))

	if status.Poll.HasEnded {
		buffer.WriteString("\nPoll has ended!\n")
	}

	voterCount := status.VoterCount
	voterUsernames := []string{}

	for _, voter := range voters {
		voterUsernames = append(voterUsernames, voter.Username)
	}

	if voterCount == 1 {
		buffer.WriteString("\n**1 user has voted**\n\n")
	} else {
		buffer.WriteString(fmt.Sprintf("\n**%d users have voted**\n\n", voterCount))
	}

	buffer.WriteString(strings.Join(voterUsernames, ", "))

	return buffer.String()
}

func FormatPollReply(poll *votingpb.Poll, reply *selectionpb.CreateSelectionReply) *contract.Embed {
	embed := &contract.Embed{
		Color: 0x5dadec,
	}

	title := ":inbox_tray: "

	if poll.Title != "" {
		title += fmt.Sprintf("**%s** ", poll.Title)
	}

	title += fmt.Sprintf("`%s`\n\n", poll.ShortId)

	embed.Title = title

	embed.Description = fmt.Sprintf("You can vote for %d option", poll.AllowedUniqueVotes)
	if poll.AllowedUniqueVotes != 1 {
		embed.Description += "s"
	}
	embed.Description += "."

	for _, batch := range reply.Batches {
		buffer := bytes.Buffer{}

		for _, batchOption := range batch.Options {
			url, hasURLMetadata := batchOption.Option.Metadata["url"]

			if hasURLMetadata {
				buffer.WriteString(fmt.Sprintf("%d. [%s](%s)\n", batchOption.Number, batchOption.Option.Content, url))
			} else {
				buffer.WriteString(fmt.Sprintf("%d. %s\n", batchOption.Number, batchOption.Option.Content))
			}
		}

		field := &contract.EmbedField{
			Name:  "\u2800",
			Value: buffer.String(),
		}

		embed.Fields = append(embed.Fields, field)
	}

	embed.Footer = &contract.EmbedFooter{
		Text: FormatVoteHelp(poll.AllowedUniqueVotes),
	}

	return embed
}

func FormatVoteHelp(allowedVotes int32) string {
	buffer := bytes.Buffer{}

	buffer.WriteString("\nTo vote type !vote followed by the numbers you would like to vote for ")

	if allowedVotes > 1 {
		buffer.WriteString("in order of your most to least favorite. e.g. \"!vote 1 2\" to vote for options 1 and 2 and you prefer option 1.")
		buffer.WriteString("\n\nOnly your most recent !vote will be counted.")
	} else {
		buffer.WriteString("e.g. !vote 1 to vote for option 1.")
	}

	return buffer.String()
}

func FormatEndPollReply(endPollReply *votingpb.EndPollReply) string {
	return fmt.Sprintf("ended poll `%s` %s", endPollReply.Poll.ShortId, endPollReply.Poll.Title)
}

func FormatVoteReply(poll *votingpb.Poll, voteReply *votingpb.VoteReply) string {
	if !voteReply.Success {
		return voteReply.Message
	}

	buffer := bytes.Buffer{}

	buffer.WriteString(":ballot_box_with_check: ")

	if poll.Title != "" {
		buffer.WriteString(fmt.Sprintf("**%s** ", poll.Title))
	}
	buffer.WriteString(fmt.Sprintf("`%s`\n\n", poll.ShortId))

	for _, voteReplyOption := range voteReply.Options {
		buffer.WriteString(fmt.Sprintf("%d. %s\n", voteReplyOption.Rank+1, voteReplyOption.Option.Content))
	}

	buffer.WriteString("\nYour vote has been submitted!")

	return buffer.String()
}

func FormatCountResult(countReply *votingpb.CountReply) *contract.Message {

	embed := &contract.Embed{
		Title:       ":ballot_box: Election Result",
		Description: fmt.Sprintf("%s `%s`", countReply.Poll.Title, countReply.Poll.ShortId),
		Color:       0x5dadec,
	}

	method := &contract.EmbedField{
		Name:  "Counting Method",
		Value: countReply.Method,
	}

	embed.Fields = append(embed.Fields, method)

	chunkedCandidates := chunkCandidates(countReply.Elected, 10)
	for _, chunk := range chunkedCandidates {
		title := "Elected"

		if len(chunkedCandidates) > 1 {
			firstRank := chunk[0].Rank
			lastRank := chunk[len(chunk)-1].Rank

			section := fmt.Sprintf(" (%d - %d)", firstRank, lastRank)

			title = title + section
		}

		buffer := bytes.Buffer{}

		for _, candidate := range chunk {
			buffer.WriteString(fmt.Sprintf("%d. %s\n", candidate.Rank, candidate.Option.Content))
		}

		field := &contract.EmbedField{
			Name:  title,
			Value: buffer.String(),
		}

		embed.Fields = append(embed.Fields, field)
	}

	message := &contract.Message{
		Embed: embed,
	}

	if len(countReply.Events) > 0 {
		events := bytes.Buffer{}

		for _, event := range countReply.Events {
			events.WriteString(event.Description + "\r\n")
		}

		fileName := fmt.Sprintf("%s_%s_results.txt", countReply.Poll.Title, countReply.Method)
		file := &contract.File{
			Name:  fileName,
			Bytes: events.Bytes(),
		}

		if len(file.Bytes) <= 2000000 {
			message.Files = append(message.Files, file)
		} else {
			message.Content = "Did not converge within a reasonable file size :/"
		}
	}

	return message
}

func FormatParseError(err error) (*contract.Response, error) {
	switch err.(type) {
	case ParseError:
		return contract.StringResponse(err.Error()), nil
	}

	return nil, err
}

func FormatClientError(err error) (*contract.Response, error) {
	st, ok := status.FromError(err)
	if !ok {
		return nil, err
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return contract.StringResponse(st.Message()), nil
	case codes.NotFound:
		return contract.StringResponse(st.Message()), nil
	}

	return nil, err
}

func chunkCandidates(candidates []*votingpb.VoteReplyOption, chunkSize int) [][]*votingpb.VoteReplyOption {
	chunked := [][]*votingpb.VoteReplyOption{}

	numCandidates := len(candidates)

	for i := 0; i < numCandidates; i += chunkSize {
		nextChunkBound := i + chunkSize

		if nextChunkBound > numCandidates {
			nextChunkBound = numCandidates
		}

		chunk := candidates[i:nextChunkBound]

		chunked = append(chunked, chunk)
	}

	return chunked
}

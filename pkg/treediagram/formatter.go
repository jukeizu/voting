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

var CountdownURL = "https://countdown.treediagram.xyz"
var BallotBoxThumbnailURL = "https://cdn.discordapp.com/attachments/320660733740449792/728375524090576996/ff85a1aae50ad48506e3275656768e89.png"
var InboxTrayThumbnailURL = "https://cdn.discordapp.com/attachments/320660733740449792/729234150829195284/d2847ce64775bce206d39ae4114db659.png"
var BallotBoxWithCheckURL = "https://cdn.discordapp.com/attachments/314238080355926016/949908877665775626/86c16c39d96283551fd4ca7392e22681_1.png"

func FormatNewPollReply(poll *votingpb.Poll) *contract.Message {
	embed := &contract.Embed{
		Color:        0x5865f2,
		Title:        "**A new poll has started**",
		ThumbnailUrl: BallotBoxThumbnailURL,
		Footer: &contract.EmbedFooter{
			Text: poll.ShortId,
		},
	}

	if poll.Title != "" {
		embed.Fields = append(embed.Fields, &contract.EmbedField{
			Name:  "Title",
			Value: poll.Title,
		})
	}

	if hasExpiration(poll) {
		embed.Fields = append(embed.Fields, &contract.EmbedField{
			Name:  "Open until",
			Value: fmt.Sprintf("<t:%d>", poll.Expires),
		})
	}

	return &contract.Message{
		Embed: embed,
		Compontents: &contract.Components{
			ActionsRows: []*contract.ActionsRow{
				&contract.ActionsRow{
					Buttons: []*contract.Button{
						&contract.Button{
							Label:    "Vote",
							CustomId: fmt.Sprintf("poll.%s", poll.ShortId),
							Emoji: contract.ComponentEmoji{
								Name: "🗳️",
							},
						},
						&contract.Button{
							Label:    "View status",
							Style:    2,
							CustomId: "pollstatus." + poll.ShortId,
						},
					},
				},
			},
		},
	}
}

func FormatPollStatusReply(status *votingpb.StatusReply, voters []*votingpb.Voter, countReply *votingpb.CountReply) *contract.Embed {
	embed := &contract.Embed{
		Color:        0x5865f2,
		ThumbnailUrl: BallotBoxThumbnailURL,
		Footer: &contract.EmbedFooter{
			Text: status.Poll.ShortId,
		},
	}

	poll := status.Poll

	if poll.Title != "" {
		embed.Fields = append(embed.Fields, &contract.EmbedField{
			Name:  "Title",
			Value: poll.Title,
		})
	}

	if status.Poll.HasEnded {
		embed.Title = "Poll has ended!"
		embed.Footer.Text = "Use !electioncount for a custom count.\n\n" + status.Poll.ShortId
		if countReply != nil {
			results := generateCountResultsEmbedField(countReply)
			embed.Fields = append(embed.Fields, results)
		}
	} else {
		embed.Title = "Poll is open!"

		if hasExpiration(poll) {
			embed.Fields = append(embed.Fields, &contract.EmbedField{
				Name:  "Open until",
				Value: fmt.Sprintf("<t:%d>", poll.Expires),
			})
		}
	}

	if status.VoterCount > 0 {
		votersField := generateVotersEmbedField(status.VoterCount, voters)
		embed.Fields = append(embed.Fields, votersField)
	}

	return embed
}

func FormatPollReply(poll *votingpb.Poll, reply *selectionpb.CreateSelectionReply) *contract.Embed {
	embed := &contract.Embed{
		Color:        0x5865f2,
		Title:        generateTitle(poll),
		ThumbnailUrl: InboxTrayThumbnailURL,
	}

	description := fmt.Sprintf("You can vote for %d option", poll.AllowedUniqueVotes)
	if poll.AllowedUniqueVotes != 1 {
		description += "s"
	}
	description += "."

	embed.Description = description

	for _, batch := range reply.Batches {
		buffer := bytes.Buffer{}

		for _, batchOption := range batch.Options {
			url, _ := batchOption.Option.Metadata["url"]
			if url != "" {
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
		Text: FormatVoteHelp(poll.AllowedUniqueVotes) + "\n\n" + poll.ShortId,
	}

	return embed
}

func FormatVoteHelp(allowedVotes int32) string {
	buffer := bytes.Buffer{}

	buffer.WriteString("\nTo vote type \"!vote\" followed by the numbers you would like to vote for")

	if allowedVotes > 1 {
		buffer.WriteString(" in order of your most to least favorite. e.g. \"!vote 1 2\" to vote for options 1 and 2 and you prefer option 1.")
		buffer.WriteString("\n\nOnly your most recent !vote will be counted.")
	} else {
		buffer.WriteString(". e.g. \"!vote 1\" to vote for option 1.")
	}

	return buffer.String()
}

func FormatOpenPollReply(openPollReply *votingpb.OpenPollReply) *contract.Embed {
	embed := &contract.Embed{
		Color:        0x5865f2,
		ThumbnailUrl: BallotBoxThumbnailURL,
		Footer: &contract.EmbedFooter{
			Text: openPollReply.Poll.ShortId,
		},
	}

	text := "Poll is open!"
	if openPollReply.PreviouslyEnded {
		text = "Poll has reopened!"
	} else if openPollReply.PreviousExpiration != openPollReply.Poll.Expires && hasExpiration(openPollReply.Poll) {
		text = "End time has changed!"
	}

	embed.Title = text

	if openPollReply.Poll.Title != "" {
		embed.Fields = append(embed.Fields, &contract.EmbedField{
			Name:  "Title",
			Value: openPollReply.Poll.Title,
		})
	}

	if hasExpiration(openPollReply.Poll) {
		embed.Fields = append(embed.Fields, &contract.EmbedField{
			Name:  "Open until",
			Value: fmt.Sprintf("<t:%d>", openPollReply.Poll.Expires),
		})
	}

	return embed
}

func FormatVoteReply(poll *votingpb.Poll, voteReply *votingpb.VoteReply) *contract.Embed {
	embed := &contract.Embed{
		Color: 0x5865f2,
		Title: ":ballot_box_with_check: Vote submitted!",
		Footer: &contract.EmbedFooter{
			Text: poll.ShortId,
		},
	}

	if !voteReply.Success {
		embed.Description = voteReply.Message
		return embed
	}

	if poll.Title != "" {
		embed.Fields = append(embed.Fields, &contract.EmbedField{
			Name:  "Title",
			Value: poll.Title,
		})
	}

	buffer := bytes.Buffer{}

	for _, voteReplyOption := range voteReply.Options {
		buffer.WriteString(fmt.Sprintf("%d. %s\n", voteReplyOption.Rank+1, voteReplyOption.Option.Content))
	}

	embed.Fields = append(embed.Fields, &contract.EmbedField{
		Name:  "Ballot",
		Value: buffer.String(),
	})

	return embed
}

func FormatCountResult(countReply *votingpb.CountReply) *contract.Message {
	embed := &contract.Embed{
		Title: ":ballot_box: Election Result",
		Color: 0x5865f2,
		Footer: &contract.EmbedFooter{
			Text: countReply.Poll.ShortId,
		},
	}

	if countReply.Poll.Title != "" {
		embed.Fields = append(embed.Fields, &contract.EmbedField{
			Name:  "Title",
			Value: countReply.Poll.Title,
		})
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

func FormatExportResult(exportReply *votingpb.ExportReply) *contract.Message {
	title := exportReply.Poll.Title
	if len(title) < 1 {
		title = exportReply.Poll.ShortId
	}

	fileName := fmt.Sprintf("%s_%s.txt", title, exportReply.Method)

	content := bytes.Buffer{}
	content.WriteString(exportReply.Content)

	file := &contract.File{
		Name:  fileName,
		Bytes: content.Bytes(),
	}

	message := &contract.Message{}

	if len(file.Bytes) <= 2000000 {
		message.Files = append(message.Files, file)
	} else {
		message.Content = "BLT is too large to upload :/"
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

func FormatClientErrorWithMessage(err error, additionalContent string) (*contract.Response, error) {
	st, ok := status.FromError(err)
	if !ok {
		return nil, err
	}

	additional := ""
	if len(additionalContent) > 0 {
		additional = "\n" + additionalContent
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return contract.StringResponse(st.Message() + additional), nil
	case codes.NotFound:
		return contract.StringResponse(st.Message() + additional), nil
	}

	return nil, err
}

func FormatClientError(err error) (*contract.Response, error) {
	return FormatClientErrorWithMessage(err, "")
}

func generateTitle(poll *votingpb.Poll) string {
	title := ""

	if poll.Title != "" {
		title += fmt.Sprintf("**%s** ", poll.Title)
	} else {
		title = "Ballot"
	}

	return title
}

func generateVotersEmbedField(voterCount int64, voters []*votingpb.Voter) *contract.EmbedField {
	voterUsernames := []string{}
	for _, voter := range voters {
		voterUsernames = append(voterUsernames, voter.Username)
	}

	votersField := &contract.EmbedField{
		Value: strings.Join(voterUsernames, ", "),
	}

	if voterCount == 1 {
		votersField.Name = "1 user voted"
	} else {
		votersField.Name = fmt.Sprintf("%d users voted", voterCount)
	}

	return votersField
}

func generateCountResultsEmbedField(countReply *votingpb.CountReply) *contract.EmbedField {
	buffer := bytes.Buffer{}

	for _, candidate := range countReply.Elected {
		buffer.WriteString(fmt.Sprintf("\n%d. %s", candidate.Rank, candidate.Option.Content))
	}

	results := &contract.EmbedField{
		Name:  "Results",
		Value: buffer.String(),
	}

	return results
}

func hasExpiration(poll *votingpb.Poll) bool {
	return poll.Expires > (time.Time{}).Unix()
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

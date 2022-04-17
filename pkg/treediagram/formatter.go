package treediagram

import (
	"bytes"
	"fmt"
	"sort"
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

	embed.Fields = append(embed.Fields, &contract.EmbedField{
		Name:  "Started by",
		Value: fmt.Sprintf("<@!%s>", poll.CreatorId),
	})

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
					},
				},
			},
		},
	}
}

func FormatPollStatusReply(status *votingpb.StatusReply, voters []*votingpb.Voter, countReply *votingpb.CountReply, editMessageId string, private bool) *contract.Message {
	embed := &contract.Embed{
		Color:        0x5865f2,
		ThumbnailUrl: BallotBoxThumbnailURL,
		Footer: &contract.EmbedFooter{
			Text: status.Poll.ShortId,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	poll := status.Poll

	if poll.Title != "" {
		embed.Fields = append(embed.Fields, &contract.EmbedField{
			Name:  "Title",
			Value: poll.Title,
		})
	}

	title := ""

	if status.Poll.HasEnded {
		title = "Poll has ended!"
		embed.Footer.Text = "Use !electioncount for a custom count.\n\n" + status.Poll.ShortId
		if countReply != nil {
			results := generateCountResultsEmbedField(countReply)
			embed.Fields = append(embed.Fields, results)
		}
	} else {
		title = "Poll is open!"

		if hasExpiration(poll) {
			embed.Fields = append(embed.Fields, &contract.EmbedField{
				Name:  "Open until",
				Value: fmt.Sprintf("<t:%d>", poll.Expires),
			})
		}
	}

	embed.Title = fmt.Sprintf("**%s**", title)

	if len(voters) > 0 {
		votersField := generateVotersEmbedField(status.VoterCount, voters)
		embed.Fields = append(embed.Fields, votersField)
	}

	message := &contract.Message{
		Embed:            embed,
		EditMessageId:    editMessageId,
		IsPrivateMessage: private,
	}

	if status.Poll.HasEnded {
		message.Compontents = &contract.Components{}
	} else {
		message.Compontents = &contract.Components{
			ActionsRows: []*contract.ActionsRow{
				&contract.ActionsRow{
					Buttons: []*contract.Button{
						&contract.Button{
							Label:    "Vote",
							CustomId: fmt.Sprintf("poll.%s", status.Poll.ShortId),
							Emoji: contract.ComponentEmoji{
								Name: "🗳️",
							},
						},
						&contract.Button{
							Label:    "Refresh",
							Style:    2,
							CustomId: fmt.Sprintf("pollstatus-refresh.%s", status.Poll.ShortId),
						},
					},
				},
			},
		}
	}

	return message
}

func FormatPollReply(poll *votingpb.Poll, reply *selectionpb.CreateSelectionReply) *contract.Embed {
	embed := &contract.Embed{
		Color:        0x5865f2,
		Title:        "Vote for " + generateTitle(poll),
		ThumbnailUrl: InboxTrayThumbnailURL,
	}

	if poll.AllowedUniqueVotes > 1 && int(poll.AllowedUniqueVotes) != len(poll.Options) {
		embed.Description = fmt.Sprintf("You can vote for %d option", poll.AllowedUniqueVotes)
	}

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

func FormatOpenPollReply(openPollReply *votingpb.OpenPollReply) *contract.Message {
	embed := &contract.Embed{
		Color:        0x5865f2,
		ThumbnailUrl: BallotBoxThumbnailURL,
		Footer: &contract.EmbedFooter{
			Text: openPollReply.Poll.ShortId,
		},
	}

	title := "Poll is open!"
	if openPollReply.PreviouslyEnded {
		title = "Poll has reopened!"
	} else if openPollReply.PreviousExpiration != openPollReply.Poll.Expires && hasExpiration(openPollReply.Poll) {
		title = "End time has changed!"
	}

	embed.Title = fmt.Sprintf("**%s**", title)

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

	embed.Fields = append(embed.Fields, &contract.EmbedField{
		Name:  "Modified by",
		Value: fmt.Sprintf("<@!%s>", openPollReply.Poll.CreatorId),
	})

	message := &contract.Message{
		Embed: embed,
		Compontents: &contract.Components{
			ActionsRows: []*contract.ActionsRow{
				&contract.ActionsRow{
					Buttons: []*contract.Button{
						&contract.Button{
							Label:    "Vote",
							CustomId: fmt.Sprintf("poll.%s", openPollReply.Poll.ShortId),
							Emoji: contract.ComponentEmoji{
								Name: "🗳️",
							},
						},
					},
				},
			},
		},
	}

	return message
}

func FormatVoteReply(poll *votingpb.Poll, voteReply *votingpb.VoteReply) *contract.Embed {
	embed := &contract.Embed{
		Color: 0x5865f2,
		Title: ":ballot_box_with_check: **Vote received!**",
		Footer: &contract.EmbedFooter{
			Text: poll.ShortId,
		},
	}

	if !voteReply.Success {
		embed.Description = voteReply.Message
		return embed
	}

	buffer := bytes.Buffer{}

	for _, voteReplyOption := range voteReply.Options {
		buffer.WriteString(fmt.Sprintf("%d. %s\n", voteReplyOption.Rank+1, voteReplyOption.Option.Content))
	}

	embed.Description = buffer.String()

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
	if poll.Title != "" {
		return poll.Title
	}

	return "Ballot"
}

func generateVotersEmbedField(voterCount int64, voters []*votingpb.Voter) *contract.EmbedField {
	voterUsernames := make([]string, len(voters))
	for i, voter := range voters {
		voterUsernames[i] = voter.Username
	}

	sort.Slice(voterUsernames, func(i, j int) bool {
		return strings.ToLower(voterUsernames[i]) < strings.ToLower(voterUsernames[j])
	})

	votersField := &contract.EmbedField{
		Value: strings.Join(voterUsernames, ", "),
	}

	votersField.Name = fmt.Sprintf("Voters: %d", voterCount)

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

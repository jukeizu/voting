package treediagram

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

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

	buffer.WriteString(fmt.Sprintf("\nType `!poll` to view the poll. A previous poll can be viewed via id. e.g. `!poll %s`", poll.ShortId))

	return buffer.String()
}

func FormatPollStatusReply(status *votingpb.StatusReply) string {
	buffer := bytes.Buffer{}

	buffer.WriteString(":ballot_box: ")
	if status.Poll.Title != "" {
		buffer.WriteString(fmt.Sprintf("**%s** ", status.Poll.Title))
	}

	buffer.WriteString(fmt.Sprintf("`%s`\n", status.Poll.ShortId))

	voterCount := len(status.Voters)

	voters := []string{}

	for _, voter := range status.Voters {
		voters = append(voters, voter.Username)
	}

	buffer.WriteString(fmt.Sprintf("\n**%d users have voted**\n\n", voterCount))
	buffer.WriteString(strings.Join(voters, ", "))

	return buffer.String()
}

func FormatPollReply(poll *votingpb.Poll, reply *selectionpb.CreateSelectionReply) string {
	buffer := bytes.Buffer{}

	options := make([]int, 0, len(reply.Options))
	for k := range reply.Options {
		options = append(options, int(k))
	}

	sort.Ints(options)

	buffer.WriteString(":inbox_tray: ")
	if poll.Title != "" {
		buffer.WriteString(fmt.Sprintf("**%s** ", poll.Title))
	}
	buffer.WriteString(fmt.Sprintf("`%s`\n\n", poll.ShortId))

	buffer.WriteString(fmt.Sprintf("You can vote for %d option(s).\n\n", poll.AllowedUniqueVotes))

	for _, k := range options {
		buffer.WriteString(fmt.Sprintf("%d. %s\n", k, reply.Options[int32(k)].Content))
	}

	buffer.WriteString(FormatVoteHelp(poll.AllowedUniqueVotes))

	return buffer.String()
}

func FormatVoteHelp(allowedVotes int32) string {
	buffer := bytes.Buffer{}

	buffer.WriteString("\nTo vote type `!vote` followed by the numbers you would like to vote for ")

	if allowedVotes > 1 {
		buffer.WriteString("**in order of your most to least favorite**. e.g. `!vote 1 2` to vote for options 1 and 2 and you prefer option 1.")
		buffer.WriteString("\n\nOnly your most recent `!vote` will be counted. Be sure to include all options you wish to vote for in a single command.")
	} else {
		buffer.WriteString("e.g. `!vote 1` to vote for option 1.")
	}

	return buffer.String()
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

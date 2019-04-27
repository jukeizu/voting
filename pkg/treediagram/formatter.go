package treediagram

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
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

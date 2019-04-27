package treediagram

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jukeizu/contract"
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

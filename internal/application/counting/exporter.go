package counting

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shawntoffel/meekstv"
)

type exporter struct {
	title  string
	config meekstv.Config
}

func (e *exporter) buildHeaderLine() string {
	return fmt.Sprintf("%d %d\n", len(e.config.Candidates), e.config.Seats)
}

func (e *exporter) buildWithdrawnCandidatesLine() string {
	if len(e.config.WithdrawnCandidates) < 1 {
		return ""
	}

	withdrawn := make([]string, len(e.config.WithdrawnCandidates))
	for i, wc := range e.config.WithdrawnCandidates {
		withdrawn[i] = strconv.Itoa(wc * -1)
	}

	return strings.Join(withdrawn, " ") + "\n"
}

func (e *exporter) buildBallotLines() string {
	sb := strings.Builder{}

	for _, ballot := range e.config.Ballots {
		fmt.Fprintf(&sb, "%d %s\n", ballot.Count, newBltBallotLine(ballot))
	}

	sb.WriteString("0\n")

	return sb.String()
}

func newBltBallotLine(ballot meekstv.Ballot) string {
	sb := strings.Builder{}

	for _, pref := range ballot.Preferences {
		fmt.Fprintf(&sb, "%d ", pref)
	}

	sb.WriteString("0")

	return sb.String()
}

func (e *exporter) buildCandidateLines() string {
	sb := strings.Builder{}

	for _, candidate := range e.config.Candidates {
		name := strings.TrimSpace(strings.ReplaceAll(candidate, "\"", ""))
		sb.WriteString("\"" + name + "\"\n")
	}

	return sb.String()
}

func (e *exporter) buildTitle() string {
	title := e.title
	if len(title) < 1 {
		title = "Election"
	}

	return "\"" + title + "\""
}

func ExportBlt(title string, config meekstv.Config) string {
	e := exporter{
		title:  title,
		config: config,
	}

	sb := strings.Builder{}

	sb.WriteString(e.buildHeaderLine())
	sb.WriteString(e.buildWithdrawnCandidatesLine())
	sb.WriteString(e.buildBallotLines())
	sb.WriteString(e.buildCandidateLines())
	sb.WriteString(e.buildTitle())

	return sb.String()
}

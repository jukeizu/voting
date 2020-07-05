package voting

import "fmt"

var (
	ErrPollHasEnded              = ValidationError{Message: "poll has ended"}
	ErrPollHasNotEnded           = ValidationError{Message: "poll has not ended"}
	ErrPastPollExpiration        = ValidationError{Message: "poll expiration must be in the future"}
	ErrNotOwner                  = ValidationError{Message: "only the poll creator may modify the poll"}
	ErrNoOptions                 = ValidationError{Message: "at least one option must be provided"}
	ErrInvalidOrDuplicateOptions = ValidationError{Message: "vote contains invalid or duplicate options"}
	ErrNoVoterExternalId         = ValidationError{Message: "voter external id required"}
)

func ErrTooManyVotes(max int32) ValidationError {
	return ValidationError{Message: fmt.Sprintf("too many votes. Maximum for this poll is %d.", max)}
}

func ErrPollNotFound(shortId string) NotFoundError {
	return NotFoundError{Message: "couldn't find poll with id: " + shortId}
}

func ErrVoterNotPermitted(voter Voter) ValidationError {
	return ValidationError{Message: fmt.Sprintf("not permitted. %s is not permitted to vote.", voter.Username)}
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

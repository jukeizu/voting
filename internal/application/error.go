package application

import "fmt"

type NotFound struct{}

func (e NotFound) Error() string {
	return "Not Found"
}

type PermissionError struct {
	Message string
}

func (e PermissionError) Error() string {
	return e.Message
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func ErrEmpty(field string) ValidationError {
	return ValidationError{
		Message: field + " is required",
	}
}

func ErrUnknownCountMethod(method string) ValidationError {
	return ValidationError{Message: fmt.Sprintf("Counting method '%s' is not supported", method)}
}

var (
	ErrPollHasEnded              = ValidationError{Message: "Poll is not open."}
	ErrPollHasNotEnded           = ValidationError{Message: "Poll has not ended."}
	ErrPastPollExpiration        = ValidationError{Message: "Poll expiration must be in the future."}
	ErrNotOwner                  = PermissionError{Message: "Only the poll creator may modify the poll."}
	ErrNoCandidates              = ValidationError{Message: "At least one option must be provided."}
	ErrInvalidOrDuplicateOptions = ValidationError{Message: "Vote contains invalid or duplicate options."}
	ErrMaxConcurrent             = ValidationError{Message: "A new poll cannot start until the others have ended."}
	ErrTooManyOptions            = ValidationError{Message: "Only one option may be provided."}
	ErrCannotVote                = PermissionError{Message: "Unable to participate."}
)

package voting

var (
	ErrPollHasEnded    = ValidationError{Message: "poll has ended"}
	ErrPollHasNotEnded = ValidationError{Message: "poll has not ended"}
	ErrNotOwner        = ValidationError{Message: "only the poll creator may end the poll"}
	ErrNoOptions       = ValidationError{Message: "at least one option must be provided"}
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

package voting

const (
	ErrPollHasEnded    = ValidationError("poll has ended")
	ErrPollHasNotEnded = ValidationError("poll has not ended")
	ErrNotOwner        = ValidationError("only the poll creator may end the poll")
)

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

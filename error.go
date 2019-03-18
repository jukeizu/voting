package voting

const (
	ErrPollHasEnded    = ValidationError("poll has ended")
	ErrPollHasNotEnded = ValidationError("poll has not ended")
)

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

package entities

type Poll struct {
	Id                 string
	Title              string
	CreatorId          string
	AllowedUniqueVotes int32
	HasEnded           bool
	Options            []Option
}

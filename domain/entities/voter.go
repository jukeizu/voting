package entities

type Voter struct {
	Id         string
	ExternalId string
	Username   string
	CanVote    bool
	Created    int32
	Updated    *int32
}

package poll

import "github.com/jukeizu/voting/api/protobuf-spec/pollpb"

type Service interface {
	Create(*pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error)
	Poll(*pollpb.PollRequest) (*pollpb.PollReply, error)
	Options(*pollpb.OptionsRequest) (*pollpb.OptionsReply, error)
	End(*pollpb.EndPollRequest) (*pollpb.EndPollReply, error)
}

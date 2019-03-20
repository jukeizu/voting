// Code generated by protoc-gen-go. DO NOT EDIT.
// source: voting.proto

package votingpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CreatePollRequest struct {
	CreatorId            string    `protobuf:"bytes,1,opt,name=creatorId,proto3" json:"creatorId,omitempty"`
	Title                string    `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	AllowedUniqueVotes   int32     `protobuf:"varint,3,opt,name=allowedUniqueVotes,proto3" json:"allowedUniqueVotes,omitempty"`
	Options              []*Option `protobuf:"bytes,4,rep,name=options,proto3" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CreatePollRequest) Reset()         { *m = CreatePollRequest{} }
func (m *CreatePollRequest) String() string { return proto.CompactTextString(m) }
func (*CreatePollRequest) ProtoMessage()    {}
func (*CreatePollRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{0}
}
func (m *CreatePollRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreatePollRequest.Unmarshal(m, b)
}
func (m *CreatePollRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreatePollRequest.Marshal(b, m, deterministic)
}
func (dst *CreatePollRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreatePollRequest.Merge(dst, src)
}
func (m *CreatePollRequest) XXX_Size() int {
	return xxx_messageInfo_CreatePollRequest.Size(m)
}
func (m *CreatePollRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreatePollRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreatePollRequest proto.InternalMessageInfo

func (m *CreatePollRequest) GetCreatorId() string {
	if m != nil {
		return m.CreatorId
	}
	return ""
}

func (m *CreatePollRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *CreatePollRequest) GetAllowedUniqueVotes() int32 {
	if m != nil {
		return m.AllowedUniqueVotes
	}
	return 0
}

func (m *CreatePollRequest) GetOptions() []*Option {
	if m != nil {
		return m.Options
	}
	return nil
}

type CreatePollReply struct {
	Poll                 *Poll    `protobuf:"bytes,1,opt,name=poll,proto3" json:"poll,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreatePollReply) Reset()         { *m = CreatePollReply{} }
func (m *CreatePollReply) String() string { return proto.CompactTextString(m) }
func (*CreatePollReply) ProtoMessage()    {}
func (*CreatePollReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{1}
}
func (m *CreatePollReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreatePollReply.Unmarshal(m, b)
}
func (m *CreatePollReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreatePollReply.Marshal(b, m, deterministic)
}
func (dst *CreatePollReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreatePollReply.Merge(dst, src)
}
func (m *CreatePollReply) XXX_Size() int {
	return xxx_messageInfo_CreatePollReply.Size(m)
}
func (m *CreatePollReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CreatePollReply.DiscardUnknown(m)
}

var xxx_messageInfo_CreatePollReply proto.InternalMessageInfo

func (m *CreatePollReply) GetPoll() *Poll {
	if m != nil {
		return m.Poll
	}
	return nil
}

type Option struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PollId               string   `protobuf:"bytes,2,opt,name=pollId,proto3" json:"pollId,omitempty"`
	Content              string   `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Option) Reset()         { *m = Option{} }
func (m *Option) String() string { return proto.CompactTextString(m) }
func (*Option) ProtoMessage()    {}
func (*Option) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{2}
}
func (m *Option) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Option.Unmarshal(m, b)
}
func (m *Option) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Option.Marshal(b, m, deterministic)
}
func (dst *Option) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Option.Merge(dst, src)
}
func (m *Option) XXX_Size() int {
	return xxx_messageInfo_Option.Size(m)
}
func (m *Option) XXX_DiscardUnknown() {
	xxx_messageInfo_Option.DiscardUnknown(m)
}

var xxx_messageInfo_Option proto.InternalMessageInfo

func (m *Option) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Option) GetPollId() string {
	if m != nil {
		return m.PollId
	}
	return ""
}

func (m *Option) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type PollRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PollRequest) Reset()         { *m = PollRequest{} }
func (m *PollRequest) String() string { return proto.CompactTextString(m) }
func (*PollRequest) ProtoMessage()    {}
func (*PollRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{3}
}
func (m *PollRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PollRequest.Unmarshal(m, b)
}
func (m *PollRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PollRequest.Marshal(b, m, deterministic)
}
func (dst *PollRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PollRequest.Merge(dst, src)
}
func (m *PollRequest) XXX_Size() int {
	return xxx_messageInfo_PollRequest.Size(m)
}
func (m *PollRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PollRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PollRequest proto.InternalMessageInfo

func (m *PollRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type PollReply struct {
	Poll                 *Poll    `protobuf:"bytes,1,opt,name=poll,proto3" json:"poll,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PollReply) Reset()         { *m = PollReply{} }
func (m *PollReply) String() string { return proto.CompactTextString(m) }
func (*PollReply) ProtoMessage()    {}
func (*PollReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{4}
}
func (m *PollReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PollReply.Unmarshal(m, b)
}
func (m *PollReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PollReply.Marshal(b, m, deterministic)
}
func (dst *PollReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PollReply.Merge(dst, src)
}
func (m *PollReply) XXX_Size() int {
	return xxx_messageInfo_PollReply.Size(m)
}
func (m *PollReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PollReply.DiscardUnknown(m)
}

var xxx_messageInfo_PollReply proto.InternalMessageInfo

func (m *PollReply) GetPoll() *Poll {
	if m != nil {
		return m.Poll
	}
	return nil
}

type Poll struct {
	Id                   string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title                string    `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	CreatorId            string    `protobuf:"bytes,3,opt,name=creatorId,proto3" json:"creatorId,omitempty"`
	AllowedUniqueVotes   int32     `protobuf:"varint,4,opt,name=allowedUniqueVotes,proto3" json:"allowedUniqueVotes,omitempty"`
	HasEnded             bool      `protobuf:"varint,5,opt,name=hasEnded,proto3" json:"hasEnded,omitempty"`
	Options              []*Option `protobuf:"bytes,6,rep,name=options,proto3" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Poll) Reset()         { *m = Poll{} }
func (m *Poll) String() string { return proto.CompactTextString(m) }
func (*Poll) ProtoMessage()    {}
func (*Poll) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{5}
}
func (m *Poll) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Poll.Unmarshal(m, b)
}
func (m *Poll) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Poll.Marshal(b, m, deterministic)
}
func (dst *Poll) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Poll.Merge(dst, src)
}
func (m *Poll) XXX_Size() int {
	return xxx_messageInfo_Poll.Size(m)
}
func (m *Poll) XXX_DiscardUnknown() {
	xxx_messageInfo_Poll.DiscardUnknown(m)
}

var xxx_messageInfo_Poll proto.InternalMessageInfo

func (m *Poll) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Poll) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Poll) GetCreatorId() string {
	if m != nil {
		return m.CreatorId
	}
	return ""
}

func (m *Poll) GetAllowedUniqueVotes() int32 {
	if m != nil {
		return m.AllowedUniqueVotes
	}
	return 0
}

func (m *Poll) GetHasEnded() bool {
	if m != nil {
		return m.HasEnded
	}
	return false
}

func (m *Poll) GetOptions() []*Option {
	if m != nil {
		return m.Options
	}
	return nil
}

type OptionsRequest struct {
	PollId               string   `protobuf:"bytes,1,opt,name=pollId,proto3" json:"pollId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OptionsRequest) Reset()         { *m = OptionsRequest{} }
func (m *OptionsRequest) String() string { return proto.CompactTextString(m) }
func (*OptionsRequest) ProtoMessage()    {}
func (*OptionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{6}
}
func (m *OptionsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OptionsRequest.Unmarshal(m, b)
}
func (m *OptionsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OptionsRequest.Marshal(b, m, deterministic)
}
func (dst *OptionsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OptionsRequest.Merge(dst, src)
}
func (m *OptionsRequest) XXX_Size() int {
	return xxx_messageInfo_OptionsRequest.Size(m)
}
func (m *OptionsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OptionsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OptionsRequest proto.InternalMessageInfo

func (m *OptionsRequest) GetPollId() string {
	if m != nil {
		return m.PollId
	}
	return ""
}

type OptionsReply struct {
	Options              []*Option `protobuf:"bytes,1,rep,name=options,proto3" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *OptionsReply) Reset()         { *m = OptionsReply{} }
func (m *OptionsReply) String() string { return proto.CompactTextString(m) }
func (*OptionsReply) ProtoMessage()    {}
func (*OptionsReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{7}
}
func (m *OptionsReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OptionsReply.Unmarshal(m, b)
}
func (m *OptionsReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OptionsReply.Marshal(b, m, deterministic)
}
func (dst *OptionsReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OptionsReply.Merge(dst, src)
}
func (m *OptionsReply) XXX_Size() int {
	return xxx_messageInfo_OptionsReply.Size(m)
}
func (m *OptionsReply) XXX_DiscardUnknown() {
	xxx_messageInfo_OptionsReply.DiscardUnknown(m)
}

var xxx_messageInfo_OptionsReply proto.InternalMessageInfo

func (m *OptionsReply) GetOptions() []*Option {
	if m != nil {
		return m.Options
	}
	return nil
}

type EndPollRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId               string   `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EndPollRequest) Reset()         { *m = EndPollRequest{} }
func (m *EndPollRequest) String() string { return proto.CompactTextString(m) }
func (*EndPollRequest) ProtoMessage()    {}
func (*EndPollRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{8}
}
func (m *EndPollRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EndPollRequest.Unmarshal(m, b)
}
func (m *EndPollRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EndPollRequest.Marshal(b, m, deterministic)
}
func (dst *EndPollRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EndPollRequest.Merge(dst, src)
}
func (m *EndPollRequest) XXX_Size() int {
	return xxx_messageInfo_EndPollRequest.Size(m)
}
func (m *EndPollRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EndPollRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EndPollRequest proto.InternalMessageInfo

func (m *EndPollRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *EndPollRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type EndPollReply struct {
	Poll                 *Poll    `protobuf:"bytes,1,opt,name=poll,proto3" json:"poll,omitempty"`
	Success              bool     `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	Reason               string   `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EndPollReply) Reset()         { *m = EndPollReply{} }
func (m *EndPollReply) String() string { return proto.CompactTextString(m) }
func (*EndPollReply) ProtoMessage()    {}
func (*EndPollReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{9}
}
func (m *EndPollReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EndPollReply.Unmarshal(m, b)
}
func (m *EndPollReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EndPollReply.Marshal(b, m, deterministic)
}
func (dst *EndPollReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EndPollReply.Merge(dst, src)
}
func (m *EndPollReply) XXX_Size() int {
	return xxx_messageInfo_EndPollReply.Size(m)
}
func (m *EndPollReply) XXX_DiscardUnknown() {
	xxx_messageInfo_EndPollReply.DiscardUnknown(m)
}

var xxx_messageInfo_EndPollReply proto.InternalMessageInfo

func (m *EndPollReply) GetPoll() *Poll {
	if m != nil {
		return m.Poll
	}
	return nil
}

func (m *EndPollReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *EndPollReply) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type StatusRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusRequest) Reset()         { *m = StatusRequest{} }
func (m *StatusRequest) String() string { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()    {}
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{10}
}
func (m *StatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRequest.Unmarshal(m, b)
}
func (m *StatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRequest.Marshal(b, m, deterministic)
}
func (dst *StatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRequest.Merge(dst, src)
}
func (m *StatusRequest) XXX_Size() int {
	return xxx_messageInfo_StatusRequest.Size(m)
}
func (m *StatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRequest proto.InternalMessageInfo

type StatusReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusReply) Reset()         { *m = StatusReply{} }
func (m *StatusReply) String() string { return proto.CompactTextString(m) }
func (*StatusReply) ProtoMessage()    {}
func (*StatusReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{11}
}
func (m *StatusReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusReply.Unmarshal(m, b)
}
func (m *StatusReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusReply.Marshal(b, m, deterministic)
}
func (dst *StatusReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusReply.Merge(dst, src)
}
func (m *StatusReply) XXX_Size() int {
	return xxx_messageInfo_StatusReply.Size(m)
}
func (m *StatusReply) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusReply.DiscardUnknown(m)
}

var xxx_messageInfo_StatusReply proto.InternalMessageInfo

type CreateBallotRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateBallotRequest) Reset()         { *m = CreateBallotRequest{} }
func (m *CreateBallotRequest) String() string { return proto.CompactTextString(m) }
func (*CreateBallotRequest) ProtoMessage()    {}
func (*CreateBallotRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{12}
}
func (m *CreateBallotRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateBallotRequest.Unmarshal(m, b)
}
func (m *CreateBallotRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateBallotRequest.Marshal(b, m, deterministic)
}
func (dst *CreateBallotRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateBallotRequest.Merge(dst, src)
}
func (m *CreateBallotRequest) XXX_Size() int {
	return xxx_messageInfo_CreateBallotRequest.Size(m)
}
func (m *CreateBallotRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateBallotRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateBallotRequest proto.InternalMessageInfo

type CreateBallotReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateBallotReply) Reset()         { *m = CreateBallotReply{} }
func (m *CreateBallotReply) String() string { return proto.CompactTextString(m) }
func (*CreateBallotReply) ProtoMessage()    {}
func (*CreateBallotReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{13}
}
func (m *CreateBallotReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateBallotReply.Unmarshal(m, b)
}
func (m *CreateBallotReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateBallotReply.Marshal(b, m, deterministic)
}
func (dst *CreateBallotReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateBallotReply.Merge(dst, src)
}
func (m *CreateBallotReply) XXX_Size() int {
	return xxx_messageInfo_CreateBallotReply.Size(m)
}
func (m *CreateBallotReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateBallotReply.DiscardUnknown(m)
}

var xxx_messageInfo_CreateBallotReply proto.InternalMessageInfo

type VoteRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VoteRequest) Reset()         { *m = VoteRequest{} }
func (m *VoteRequest) String() string { return proto.CompactTextString(m) }
func (*VoteRequest) ProtoMessage()    {}
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{14}
}
func (m *VoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteRequest.Unmarshal(m, b)
}
func (m *VoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteRequest.Marshal(b, m, deterministic)
}
func (dst *VoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteRequest.Merge(dst, src)
}
func (m *VoteRequest) XXX_Size() int {
	return xxx_messageInfo_VoteRequest.Size(m)
}
func (m *VoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VoteRequest proto.InternalMessageInfo

type VoteReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VoteReply) Reset()         { *m = VoteReply{} }
func (m *VoteReply) String() string { return proto.CompactTextString(m) }
func (*VoteReply) ProtoMessage()    {}
func (*VoteReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{15}
}
func (m *VoteReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteReply.Unmarshal(m, b)
}
func (m *VoteReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteReply.Marshal(b, m, deterministic)
}
func (dst *VoteReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteReply.Merge(dst, src)
}
func (m *VoteReply) XXX_Size() int {
	return xxx_messageInfo_VoteReply.Size(m)
}
func (m *VoteReply) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteReply.DiscardUnknown(m)
}

var xxx_messageInfo_VoteReply proto.InternalMessageInfo

type CountRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CountRequest) Reset()         { *m = CountRequest{} }
func (m *CountRequest) String() string { return proto.CompactTextString(m) }
func (*CountRequest) ProtoMessage()    {}
func (*CountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{16}
}
func (m *CountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CountRequest.Unmarshal(m, b)
}
func (m *CountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CountRequest.Marshal(b, m, deterministic)
}
func (dst *CountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CountRequest.Merge(dst, src)
}
func (m *CountRequest) XXX_Size() int {
	return xxx_messageInfo_CountRequest.Size(m)
}
func (m *CountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CountRequest proto.InternalMessageInfo

type CountReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CountReply) Reset()         { *m = CountReply{} }
func (m *CountReply) String() string { return proto.CompactTextString(m) }
func (*CountReply) ProtoMessage()    {}
func (*CountReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_voting_e436c270bf504a8c, []int{17}
}
func (m *CountReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CountReply.Unmarshal(m, b)
}
func (m *CountReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CountReply.Marshal(b, m, deterministic)
}
func (dst *CountReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CountReply.Merge(dst, src)
}
func (m *CountReply) XXX_Size() int {
	return xxx_messageInfo_CountReply.Size(m)
}
func (m *CountReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CountReply.DiscardUnknown(m)
}

var xxx_messageInfo_CountReply proto.InternalMessageInfo

func init() {
	proto.RegisterType((*CreatePollRequest)(nil), "votingpb.CreatePollRequest")
	proto.RegisterType((*CreatePollReply)(nil), "votingpb.CreatePollReply")
	proto.RegisterType((*Option)(nil), "votingpb.Option")
	proto.RegisterType((*PollRequest)(nil), "votingpb.PollRequest")
	proto.RegisterType((*PollReply)(nil), "votingpb.PollReply")
	proto.RegisterType((*Poll)(nil), "votingpb.Poll")
	proto.RegisterType((*OptionsRequest)(nil), "votingpb.OptionsRequest")
	proto.RegisterType((*OptionsReply)(nil), "votingpb.OptionsReply")
	proto.RegisterType((*EndPollRequest)(nil), "votingpb.EndPollRequest")
	proto.RegisterType((*EndPollReply)(nil), "votingpb.EndPollReply")
	proto.RegisterType((*StatusRequest)(nil), "votingpb.StatusRequest")
	proto.RegisterType((*StatusReply)(nil), "votingpb.StatusReply")
	proto.RegisterType((*CreateBallotRequest)(nil), "votingpb.CreateBallotRequest")
	proto.RegisterType((*CreateBallotReply)(nil), "votingpb.CreateBallotReply")
	proto.RegisterType((*VoteRequest)(nil), "votingpb.VoteRequest")
	proto.RegisterType((*VoteReply)(nil), "votingpb.VoteReply")
	proto.RegisterType((*CountRequest)(nil), "votingpb.CountRequest")
	proto.RegisterType((*CountReply)(nil), "votingpb.CountReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// VotingClient is the client API for Voting service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VotingClient interface {
	CreatePoll(ctx context.Context, in *CreatePollRequest, opts ...grpc.CallOption) (*CreatePollReply, error)
	Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollReply, error)
	EndPoll(ctx context.Context, in *EndPollRequest, opts ...grpc.CallOption) (*EndPollReply, error)
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error)
	CreateBallot(ctx context.Context, in *CreateBallotRequest, opts ...grpc.CallOption) (*CreateBallotReply, error)
	Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteReply, error)
	Count(ctx context.Context, in *CountRequest, opts ...grpc.CallOption) (*CountReply, error)
}

type votingClient struct {
	cc *grpc.ClientConn
}

func NewVotingClient(cc *grpc.ClientConn) VotingClient {
	return &votingClient{cc}
}

func (c *votingClient) CreatePoll(ctx context.Context, in *CreatePollRequest, opts ...grpc.CallOption) (*CreatePollReply, error) {
	out := new(CreatePollReply)
	err := c.cc.Invoke(ctx, "/votingpb.Voting/CreatePoll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingClient) Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollReply, error) {
	out := new(PollReply)
	err := c.cc.Invoke(ctx, "/votingpb.Voting/Poll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingClient) EndPoll(ctx context.Context, in *EndPollRequest, opts ...grpc.CallOption) (*EndPollReply, error) {
	out := new(EndPollReply)
	err := c.cc.Invoke(ctx, "/votingpb.Voting/EndPoll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error) {
	out := new(StatusReply)
	err := c.cc.Invoke(ctx, "/votingpb.Voting/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingClient) CreateBallot(ctx context.Context, in *CreateBallotRequest, opts ...grpc.CallOption) (*CreateBallotReply, error) {
	out := new(CreateBallotReply)
	err := c.cc.Invoke(ctx, "/votingpb.Voting/CreateBallot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingClient) Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteReply, error) {
	out := new(VoteReply)
	err := c.cc.Invoke(ctx, "/votingpb.Voting/Vote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingClient) Count(ctx context.Context, in *CountRequest, opts ...grpc.CallOption) (*CountReply, error) {
	out := new(CountReply)
	err := c.cc.Invoke(ctx, "/votingpb.Voting/Count", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VotingServer is the server API for Voting service.
type VotingServer interface {
	CreatePoll(context.Context, *CreatePollRequest) (*CreatePollReply, error)
	Poll(context.Context, *PollRequest) (*PollReply, error)
	EndPoll(context.Context, *EndPollRequest) (*EndPollReply, error)
	Status(context.Context, *StatusRequest) (*StatusReply, error)
	CreateBallot(context.Context, *CreateBallotRequest) (*CreateBallotReply, error)
	Vote(context.Context, *VoteRequest) (*VoteReply, error)
	Count(context.Context, *CountRequest) (*CountReply, error)
}

func RegisterVotingServer(s *grpc.Server, srv VotingServer) {
	s.RegisterService(&_Voting_serviceDesc, srv)
}

func _Voting_CreatePoll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServer).CreatePoll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/votingpb.Voting/CreatePoll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServer).CreatePoll(ctx, req.(*CreatePollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Voting_Poll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServer).Poll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/votingpb.Voting/Poll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServer).Poll(ctx, req.(*PollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Voting_EndPoll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndPollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServer).EndPoll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/votingpb.Voting/EndPoll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServer).EndPoll(ctx, req.(*EndPollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Voting_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/votingpb.Voting/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Voting_CreateBallot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBallotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServer).CreateBallot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/votingpb.Voting/CreateBallot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServer).CreateBallot(ctx, req.(*CreateBallotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Voting_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/votingpb.Voting/Vote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServer).Vote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Voting_Count_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServer).Count(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/votingpb.Voting/Count",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServer).Count(ctx, req.(*CountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Voting_serviceDesc = grpc.ServiceDesc{
	ServiceName: "votingpb.Voting",
	HandlerType: (*VotingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePoll",
			Handler:    _Voting_CreatePoll_Handler,
		},
		{
			MethodName: "Poll",
			Handler:    _Voting_Poll_Handler,
		},
		{
			MethodName: "EndPoll",
			Handler:    _Voting_EndPoll_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Voting_Status_Handler,
		},
		{
			MethodName: "CreateBallot",
			Handler:    _Voting_CreateBallot_Handler,
		},
		{
			MethodName: "Vote",
			Handler:    _Voting_Vote_Handler,
		},
		{
			MethodName: "Count",
			Handler:    _Voting_Count_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "voting.proto",
}

func init() { proto.RegisterFile("voting.proto", fileDescriptor_voting_e436c270bf504a8c) }

var fileDescriptor_voting_e436c270bf504a8c = []byte{
	// 543 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x5d, 0xd6, 0x36, 0x6d, 0x6f, 0xd2, 0x0e, 0xdc, 0x75, 0x84, 0x8c, 0x49, 0x95, 0x9f, 0x2a,
	0x1e, 0x8a, 0x54, 0x40, 0xa0, 0x49, 0xbc, 0x30, 0x0d, 0x69, 0x08, 0x09, 0x14, 0xc4, 0xde, 0xb3,
	0xc6, 0x82, 0x48, 0x56, 0x9c, 0xc5, 0x0e, 0xa8, 0xdf, 0xc3, 0x7f, 0xf0, 0xc0, 0x97, 0x21, 0xdb,
	0x71, 0xe2, 0xa6, 0x19, 0x83, 0xc7, 0x73, 0xef, 0xb9, 0xf6, 0x3d, 0xe7, 0xc4, 0x01, 0xff, 0x3b,
	0x13, 0x69, 0xf6, 0x75, 0x95, 0x17, 0x4c, 0x30, 0x34, 0xd2, 0x28, 0xbf, 0xc1, 0x3f, 0x1d, 0x78,
	0x78, 0x51, 0x90, 0x58, 0x90, 0x4f, 0x8c, 0xd2, 0x88, 0xdc, 0x96, 0x84, 0x0b, 0xf4, 0x04, 0xc6,
	0x1b, 0x59, 0x64, 0xc5, 0x55, 0x12, 0x38, 0x0b, 0x67, 0x39, 0x8e, 0x9a, 0x02, 0x3a, 0x86, 0x81,
	0x48, 0x05, 0x25, 0xc1, 0xa1, 0xea, 0x68, 0x80, 0x56, 0x80, 0x62, 0x4a, 0xd9, 0x0f, 0x92, 0x7c,
	0xc9, 0xd2, 0xdb, 0x92, 0x5c, 0x33, 0x41, 0x78, 0xd0, 0x5b, 0x38, 0xcb, 0x41, 0xd4, 0xd1, 0x41,
	0x4f, 0x61, 0xc8, 0x72, 0x91, 0xb2, 0x8c, 0x07, 0xfd, 0x45, 0x6f, 0xe9, 0xad, 0x1f, 0xac, 0xcc,
	0x56, 0xab, 0x8f, 0xaa, 0x11, 0x19, 0x02, 0x7e, 0x09, 0x47, 0xf6, 0x92, 0x39, 0xdd, 0x22, 0x0c,
	0xfd, 0x9c, 0x51, 0xaa, 0xb6, 0xf3, 0xd6, 0xd3, 0x66, 0x56, 0x51, 0x54, 0x0f, 0xbf, 0x07, 0x57,
	0x9f, 0x84, 0xa6, 0x70, 0x98, 0x1a, 0x25, 0x87, 0x69, 0x82, 0x4e, 0xc0, 0x95, 0x8c, 0xab, 0xa4,
	0xd2, 0x50, 0x21, 0x14, 0xc0, 0x70, 0xc3, 0x32, 0x41, 0x32, 0xa1, 0x36, 0x1f, 0x47, 0x06, 0xe2,
	0x33, 0xf0, 0x6c, 0x87, 0x5a, 0x07, 0xe2, 0x67, 0x30, 0xfe, 0xbf, 0xdd, 0x7e, 0x3b, 0xd0, 0x97,
	0x70, 0x6f, 0xb5, 0x6e, 0x77, 0x77, 0x12, 0xe9, 0xb5, 0x13, 0xe9, 0xf6, 0xbe, 0x7f, 0xa7, 0xf7,
	0x21, 0x8c, 0xbe, 0xc5, 0xfc, 0x32, 0x4b, 0x48, 0x12, 0x0c, 0x16, 0xce, 0x72, 0x14, 0xd5, 0xd8,
	0xce, 0xc5, 0xbd, 0x2f, 0x97, 0x25, 0x4c, 0x75, 0x89, 0x1b, 0x5f, 0x1a, 0x63, 0x1d, 0xdb, 0x58,
	0x7c, 0x0e, 0x7e, 0xcd, 0x94, 0x16, 0x59, 0xb7, 0x38, 0xf7, 0xdd, 0xf2, 0x1a, 0xa6, 0x97, 0x59,
	0xf2, 0x17, 0xf7, 0xe5, 0xad, 0x25, 0x27, 0x45, 0x13, 0xa7, 0x46, 0x38, 0x01, 0xbf, 0x9e, 0xfc,
	0xc7, 0x60, 0xe4, 0x27, 0xc0, 0xcb, 0xcd, 0x86, 0x70, 0xae, 0x0e, 0x1b, 0x45, 0x06, 0xca, 0x5b,
	0x0a, 0x12, 0x73, 0x96, 0x55, 0x01, 0x54, 0x08, 0x1f, 0xc1, 0xe4, 0xb3, 0x88, 0x45, 0x69, 0x4c,
	0xc0, 0x13, 0xf0, 0x4c, 0x21, 0xa7, 0x5b, 0x3c, 0x87, 0x99, 0xfe, 0x7a, 0xdf, 0xca, 0x28, 0x84,
	0x61, 0xcd, 0xcc, 0xcb, 0x33, 0x65, 0xc9, 0x9d, 0x80, 0x27, 0x23, 0x32, 0x1c, 0x0f, 0xc6, 0x1a,
	0xca, 0xde, 0x14, 0xfc, 0x0b, 0x56, 0x66, 0xf5, 0x01, 0x3e, 0x40, 0x85, 0x73, 0xba, 0x5d, 0xff,
	0xea, 0x81, 0x7b, 0xad, 0xf4, 0xa0, 0x77, 0x00, 0xcd, 0x73, 0x41, 0xa7, 0x8d, 0xcc, 0xbd, 0x97,
	0x1e, 0x3e, 0xee, 0x6e, 0xca, 0xeb, 0x0e, 0xd0, 0x8b, 0xea, 0x13, 0x9d, 0xb7, 0x8c, 0xaa, 0x66,
	0x67, 0xed, 0xb2, 0x9e, 0x7a, 0x03, 0xc3, 0xca, 0x74, 0x14, 0x34, 0x8c, 0xdd, 0x04, 0xc3, 0x93,
	0x8e, 0x8e, 0x1e, 0x3f, 0x07, 0x57, 0x9b, 0x87, 0x1e, 0x35, 0x9c, 0x1d, 0x7f, 0xc3, 0xf9, 0x7e,
	0x43, 0xcf, 0x7e, 0x00, 0xdf, 0xb6, 0x14, 0x9d, 0xb5, 0xd5, 0xed, 0x24, 0x10, 0x9e, 0xde, 0xd5,
	0xae, 0xe5, 0x4b, 0xf3, 0x6d, 0xf9, 0x56, 0x36, 0xb6, 0xfc, 0x26, 0xa3, 0x03, 0xf4, 0x0a, 0x06,
	0x2a, 0x15, 0x64, 0x49, 0xb4, 0x63, 0x0b, 0x8f, 0xf7, 0xea, 0x6a, 0xf0, 0xc6, 0x55, 0xff, 0xe6,
	0xe7, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x1d, 0x3a, 0x2a, 0x1d, 0xab, 0x05, 0x00, 0x00,
}
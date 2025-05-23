// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: voting.proto

package votingpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	VotingService_CreatePoll_FullMethodName = "/voting.v1.VotingService/CreatePoll"
	VotingService_Poll_FullMethodName       = "/voting.v1.VotingService/Poll"
	VotingService_OpenPoll_FullMethodName   = "/voting.v1.VotingService/OpenPoll"
	VotingService_EndPoll_FullMethodName    = "/voting.v1.VotingService/EndPoll"
	VotingService_Choices_FullMethodName    = "/voting.v1.VotingService/Choices"
	VotingService_Vote_FullMethodName       = "/voting.v1.VotingService/Vote"
	VotingService_Ballot_FullMethodName     = "/voting.v1.VotingService/Ballot"
	VotingService_Count_FullMethodName      = "/voting.v1.VotingService/Count"
	VotingService_Export_FullMethodName     = "/voting.v1.VotingService/Export"
)

// VotingServiceClient is the client API for VotingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VotingServiceClient interface {
	CreatePoll(ctx context.Context, in *CreatePollRequest, opts ...grpc.CallOption) (*CreatePollResponse, error)
	Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollResponse, error)
	OpenPoll(ctx context.Context, in *OpenPollRequest, opts ...grpc.CallOption) (*OpenPollResponse, error)
	EndPoll(ctx context.Context, in *EndPollRequest, opts ...grpc.CallOption) (*EndPollResponse, error)
	Choices(ctx context.Context, in *ChoicesRequest, opts ...grpc.CallOption) (*ChoicesResponse, error)
	Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error)
	Ballot(ctx context.Context, in *BallotRequest, opts ...grpc.CallOption) (*BallotResponse, error)
	Count(ctx context.Context, in *CountRequest, opts ...grpc.CallOption) (*CountResponse, error)
	Export(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error)
}

type votingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVotingServiceClient(cc grpc.ClientConnInterface) VotingServiceClient {
	return &votingServiceClient{cc}
}

func (c *votingServiceClient) CreatePoll(ctx context.Context, in *CreatePollRequest, opts ...grpc.CallOption) (*CreatePollResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreatePollResponse)
	err := c.cc.Invoke(ctx, VotingService_CreatePoll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingServiceClient) Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PollResponse)
	err := c.cc.Invoke(ctx, VotingService_Poll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingServiceClient) OpenPoll(ctx context.Context, in *OpenPollRequest, opts ...grpc.CallOption) (*OpenPollResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpenPollResponse)
	err := c.cc.Invoke(ctx, VotingService_OpenPoll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingServiceClient) EndPoll(ctx context.Context, in *EndPollRequest, opts ...grpc.CallOption) (*EndPollResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EndPollResponse)
	err := c.cc.Invoke(ctx, VotingService_EndPoll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingServiceClient) Choices(ctx context.Context, in *ChoicesRequest, opts ...grpc.CallOption) (*ChoicesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ChoicesResponse)
	err := c.cc.Invoke(ctx, VotingService_Choices_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingServiceClient) Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VoteResponse)
	err := c.cc.Invoke(ctx, VotingService_Vote_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingServiceClient) Ballot(ctx context.Context, in *BallotRequest, opts ...grpc.CallOption) (*BallotResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BallotResponse)
	err := c.cc.Invoke(ctx, VotingService_Ballot_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingServiceClient) Count(ctx context.Context, in *CountRequest, opts ...grpc.CallOption) (*CountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CountResponse)
	err := c.cc.Invoke(ctx, VotingService_Count_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *votingServiceClient) Export(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExportResponse)
	err := c.cc.Invoke(ctx, VotingService_Export_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VotingServiceServer is the server API for VotingService service.
// All implementations must embed UnimplementedVotingServiceServer
// for forward compatibility.
type VotingServiceServer interface {
	CreatePoll(context.Context, *CreatePollRequest) (*CreatePollResponse, error)
	Poll(context.Context, *PollRequest) (*PollResponse, error)
	OpenPoll(context.Context, *OpenPollRequest) (*OpenPollResponse, error)
	EndPoll(context.Context, *EndPollRequest) (*EndPollResponse, error)
	Choices(context.Context, *ChoicesRequest) (*ChoicesResponse, error)
	Vote(context.Context, *VoteRequest) (*VoteResponse, error)
	Ballot(context.Context, *BallotRequest) (*BallotResponse, error)
	Count(context.Context, *CountRequest) (*CountResponse, error)
	Export(context.Context, *ExportRequest) (*ExportResponse, error)
	mustEmbedUnimplementedVotingServiceServer()
}

// UnimplementedVotingServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedVotingServiceServer struct{}

func (UnimplementedVotingServiceServer) CreatePoll(context.Context, *CreatePollRequest) (*CreatePollResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePoll not implemented")
}
func (UnimplementedVotingServiceServer) Poll(context.Context, *PollRequest) (*PollResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Poll not implemented")
}
func (UnimplementedVotingServiceServer) OpenPoll(context.Context, *OpenPollRequest) (*OpenPollResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenPoll not implemented")
}
func (UnimplementedVotingServiceServer) EndPoll(context.Context, *EndPollRequest) (*EndPollResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EndPoll not implemented")
}
func (UnimplementedVotingServiceServer) Choices(context.Context, *ChoicesRequest) (*ChoicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Choices not implemented")
}
func (UnimplementedVotingServiceServer) Vote(context.Context, *VoteRequest) (*VoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
}
func (UnimplementedVotingServiceServer) Ballot(context.Context, *BallotRequest) (*BallotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ballot not implemented")
}
func (UnimplementedVotingServiceServer) Count(context.Context, *CountRequest) (*CountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Count not implemented")
}
func (UnimplementedVotingServiceServer) Export(context.Context, *ExportRequest) (*ExportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Export not implemented")
}
func (UnimplementedVotingServiceServer) mustEmbedUnimplementedVotingServiceServer() {}
func (UnimplementedVotingServiceServer) testEmbeddedByValue()                       {}

// UnsafeVotingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VotingServiceServer will
// result in compilation errors.
type UnsafeVotingServiceServer interface {
	mustEmbedUnimplementedVotingServiceServer()
}

func RegisterVotingServiceServer(s grpc.ServiceRegistrar, srv VotingServiceServer) {
	// If the following call pancis, it indicates UnimplementedVotingServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&VotingService_ServiceDesc, srv)
}

func _VotingService_CreatePoll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServiceServer).CreatePoll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VotingService_CreatePoll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServiceServer).CreatePoll(ctx, req.(*CreatePollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_Poll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServiceServer).Poll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VotingService_Poll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServiceServer).Poll(ctx, req.(*PollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_OpenPoll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenPollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServiceServer).OpenPoll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VotingService_OpenPoll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServiceServer).OpenPoll(ctx, req.(*OpenPollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_EndPoll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndPollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServiceServer).EndPoll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VotingService_EndPoll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServiceServer).EndPoll(ctx, req.(*EndPollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_Choices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChoicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServiceServer).Choices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VotingService_Choices_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServiceServer).Choices(ctx, req.(*ChoicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServiceServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VotingService_Vote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServiceServer).Vote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_Ballot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BallotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServiceServer).Ballot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VotingService_Ballot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServiceServer).Ballot(ctx, req.(*BallotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_Count_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServiceServer).Count(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VotingService_Count_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServiceServer).Count(ctx, req.(*CountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VotingService_Export_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VotingServiceServer).Export(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VotingService_Export_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VotingServiceServer).Export(ctx, req.(*ExportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VotingService_ServiceDesc is the grpc.ServiceDesc for VotingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VotingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "voting.v1.VotingService",
	HandlerType: (*VotingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePoll",
			Handler:    _VotingService_CreatePoll_Handler,
		},
		{
			MethodName: "Poll",
			Handler:    _VotingService_Poll_Handler,
		},
		{
			MethodName: "OpenPoll",
			Handler:    _VotingService_OpenPoll_Handler,
		},
		{
			MethodName: "EndPoll",
			Handler:    _VotingService_EndPoll_Handler,
		},
		{
			MethodName: "Choices",
			Handler:    _VotingService_Choices_Handler,
		},
		{
			MethodName: "Vote",
			Handler:    _VotingService_Vote_Handler,
		},
		{
			MethodName: "Ballot",
			Handler:    _VotingService_Ballot_Handler,
		},
		{
			MethodName: "Count",
			Handler:    _VotingService_Count_Handler,
		},
		{
			MethodName: "Export",
			Handler:    _VotingService_Export_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "voting.proto",
}

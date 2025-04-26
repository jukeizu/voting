package infrastructure

import (
	"context"
	"net"
	"time"

	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	"github.com/jukeizu/voting/internal/application"
	"github.com/jukeizu/voting/internal/application/counting"
	"github.com/jukeizu/voting/internal/application/polling"
	"github.com/jukeizu/voting/internal/application/registration"
	"github.com/jukeizu/voting/internal/application/voting"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type CorrelationIDKey struct{}

var _ votingpb.VotingServiceServer = &GrpcServer{}

type GrpcServer struct {
	votingpb.VotingServiceServer
	logger              zerolog.Logger
	server              *grpc.Server
	RegistrationHandler registration.Handler
	PollingHandler      polling.Handler
	VotingHandler       voting.Handler
	CountingHandler     counting.Handler
}

func NewGrpcServer(
	logger zerolog.Logger,
	registrationHandler registration.Handler,
	pollingHandler polling.Handler,
	votingHandler voting.Handler,
	countingHandler counting.Handler,
) GrpcServer {
	g := GrpcServer{
		logger:              logger,
		server:              newGrpcServer(logger, registrationHandler),
		RegistrationHandler: registrationHandler,
		PollingHandler:      pollingHandler,
		VotingHandler:       votingHandler,
		CountingHandler:     countingHandler,
	}

	votingpb.RegisterVotingServiceServer(g.server, g)
	reflection.Register(g.server)

	return g
}

func (s GrpcServer) CreatePoll(ctx context.Context, req *votingpb.CreatePollRequest) (*votingpb.CreatePollResponse, error) {
	poll, err := s.PollingHandler.CreatePoll(toCreatePollRequest(s.getVoterFromContext(ctx), req))
	if err != nil {
		return nil, toStatusErr(err)
	}

	return &votingpb.CreatePollResponse{
		Poll: toPbPoll(poll),
	}, nil
}

func (s GrpcServer) Poll(ctx context.Context, req *votingpb.PollRequest) (*votingpb.PollResponse, error) {
	poll, err := s.PollingHandler.Poll(polling.PollRequest{
		ID:    req.Id,
		Voter: s.getVoterFromContext(ctx),
	})
	if err != nil {
		return nil, toStatusErr(err)
	}

	return &votingpb.PollResponse{
		Poll: toPbPoll(poll),
	}, nil
}

func (s GrpcServer) OpenPoll(ctx context.Context, req *votingpb.OpenPollRequest) (*votingpb.OpenPollResponse, error) {
	var expires *time.Time
	if req.Expires > 0 {
		t := time.Unix(req.Expires, 0)
		expires = &t
	}

	poll, err := s.PollingHandler.Open(polling.OpenPollRequest{
		PollId:  req.Id,
		Expires: expires,
		Voter:   s.getVoterFromContext(ctx),
	})
	if err != nil {
		return nil, toStatusErr(err)
	}
	return &votingpb.OpenPollResponse{
		Poll: toPbPoll(poll),
	}, nil
}

func (s GrpcServer) EndPoll(ctx context.Context, req *votingpb.EndPollRequest) (*votingpb.EndPollResponse, error) {
	poll, err := s.PollingHandler.End(polling.EndPollRequest{
		PollId: req.Id,
		Voter:  s.getVoterFromContext(ctx),
	})
	if err != nil {
		return nil, toStatusErr(err)
	}
	return &votingpb.EndPollResponse{
		Poll: toPbPoll(poll),
	}, nil
}

func (s GrpcServer) Choices(ctx context.Context, req *votingpb.ChoicesRequest) (*votingpb.ChoicesResponse, error) {
	resp, err := s.VotingHandler.Choices(voting.ChoicesRequest{
		Voter:  s.getVoterFromContext(ctx),
		PollID: req.PollId,
	})
	if err != nil {
		return nil, toStatusErr(err)
	}

	return &votingpb.ChoicesResponse{
		Title:   resp.Title,
		Choices: toPbChoices(resp.Choices),
	}, nil
}

func (s GrpcServer) Vote(ctx context.Context, req *votingpb.VoteRequest) (*votingpb.VoteResponse, error) {
	resp, err := s.VotingHandler.SubmitVote(voting.SubmitVoteRequest{
		Voter:   s.getVoterFromContext(ctx),
		Options: toBallotOptions(req.Options),
	})
	if err != nil {
		return nil, toStatusErr(err)
	}

	return &votingpb.VoteResponse{
		Choices: toPbRankedChoices(resp.RankedChoices),
	}, nil
}

func (s GrpcServer) Ballot(ctx context.Context, req *votingpb.BallotRequest) (*votingpb.BallotResponse, error) {
	resp, err := s.VotingHandler.Ballot(voting.BallotRequest{
		Voter:  s.getVoterFromContext(ctx),
		PollID: req.PollId,
	})
	if err != nil {
		return nil, toStatusErr(err)
	}

	return &votingpb.BallotResponse{
		RankedChoices: toPbRankedChoices(resp.RankedChoices),
	}, nil
}

func (s GrpcServer) Count(ctx context.Context, req *votingpb.CountRequest) (*votingpb.CountResponse, error) {
	resp, err := s.CountingHandler.Count(counting.CountRequest{
		PollID: req.PollId,
		Seats:  int(req.Seats),
		Method: "meekstv",
	})
	if err != nil {
		return nil, toStatusErr(err)
	}
	return &votingpb.CountResponse{
		Title:      resp.Title,
		Candidates: int32(resp.Candidates),
		Ballots:    int32(resp.Ballots),
		Method:     resp.Method,
		Elected:    toPbElectedCandidates(resp.Elected),
		Report:     resp.Report,
	}, nil
}

func (s GrpcServer) Export(ctx context.Context, req *votingpb.ExportRequest) (*votingpb.ExportResponse, error) {
	withdrawn := make([]int, len(req.Withdrawn))
	for i, w := range req.Withdrawn {
		withdrawn[i] = int(w)
	}

	resp, err := s.CountingHandler.Export(counting.ExportRequest{
		PollID:    req.PollId,
		Seats:     int(req.Seats),
		Withdrawn: withdrawn,
	})
	if err != nil {
		return nil, toStatusErr(err)
	}
	return &votingpb.ExportResponse{
		Blt: resp.BLT,
	}, nil
}

func toCreatePollRequest(voter application.Voter, req *votingpb.CreatePollRequest) polling.CreatePollRequest {
	poll := polling.CreatePollRequest{
		Title:   req.Title,
		Ranked:  req.Ranked,
		Creator: voter,
	}

	if req.Expires > 0 {
		t := time.Unix(req.Expires, 0)
		poll.Expiration = &t
	}

	for _, createOption := range req.Candidates {
		poll.Candidates = append(poll.Candidates, application.Candidate{
			Name: createOption.Name,
			URL:  createOption.Url,
		})
	}

	return poll
}

func toPbPoll(poll application.Poll) *votingpb.Poll {
	pbPoll := &votingpb.Poll{
		Id:         poll.Id,
		ExternalId: poll.ExternalID(),
		Title:      poll.Title,
		HasEnded:   poll.Ended,
		Organization: &votingpb.Organization{
			Name:       poll.Organization.Name,
			ExternalId: poll.Organization.ExternalId,
		},
		Creator: &votingpb.Voter{
			Name:       poll.Creator.Name,
			ExternalId: poll.Creator.ExternalId,
		},
		Voters: toPbVoters(poll.Voters),
	}

	if poll.Expires != nil {
		pbPoll.Expires = poll.Expires.Unix()
	}

	return pbPoll
}

func toPbElectedCandidates(options []application.ElectedCandidate) []*votingpb.ElectedCandidate {
	pbOptions := []*votingpb.ElectedCandidate{}

	for _, o := range options {
		pbOptions = append(pbOptions, &votingpb.ElectedCandidate{
			Rank: o.Rank,
			Name: o.Name,
		})
	}

	return pbOptions
}

func toPbVoters(voters []application.Voter) []*votingpb.Voter {
	pbOptions := []*votingpb.Voter{}

	for _, voter := range voters {
		pbOptions = append(pbOptions, &votingpb.Voter{
			ExternalId: voter.ExternalId,
			Name:       voter.Name,
		})
	}

	return pbOptions
}

func toPbChoices(choices []application.Choice) []*votingpb.Choice {
	pbChoices := []*votingpb.Choice{}

	for _, choice := range choices {
		pbChoices = append(pbChoices, &votingpb.Choice{
			Number:  choice.Number,
			Content: choice.Name,
			Url:     choice.URL,
		})
	}

	return pbChoices
}

func toPbRankedChoices(choices []application.RankedChoice) []*votingpb.RankedChoice {
	pbRankedChoices := []*votingpb.RankedChoice{}
	for _, rankedChoice := range choices {
		pbRankedChoices = append(pbRankedChoices, &votingpb.RankedChoice{
			Rank: rankedChoice.Rank,
			Choice: &votingpb.Choice{
				Number:  rankedChoice.Choice.Number,
				Content: rankedChoice.Choice.Name,
				Url:     rankedChoice.Choice.URL,
			},
		})
	}
	return pbRankedChoices
}

func toBallotOptions(pbOptions []*votingpb.BallotOption) []application.BallotOption {
	ballotOptions := make([]application.BallotOption, len(pbOptions))

	for i, pbOption := range pbOptions {
		ballotOption := application.BallotOption{
			Number: pbOption.Number,
			Rank:   pbOption.Rank,
		}

		ballotOptions[i] = ballotOption
	}

	return ballotOptions
}

func toStatusErr(err error) error {
	switch err.(type) {
	case application.ValidationError:
		return status.Error(codes.InvalidArgument, err.Error())
	case application.PermissionError:
		return status.Error(codes.PermissionDenied, err.Error())
	case application.NotFound:
		return status.Error(codes.NotFound, err.Error())
	}
	return err
}

func (s GrpcServer) getVoterFromContext(ctx context.Context) application.Voter {
	return ctx.Value(application.Voter{}).(application.Voter)
}

func (s GrpcServer) Start(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer s.Stop()

	errors := make(chan error, 1)

	go func() {
		err := s.server.Serve(listener)
		if err != nil {
			errors <- err
		}
	}()

	s.logger.Info().
		Str("transport", "grpc").
		Str("addr", addr).
		Msg("listening")

	select {
	case <-ctx.Done():
		return nil
	case err := <-errors:
		return err
	}
}

func (s GrpcServer) Stop() {
	if s.server == nil {
		return
	}

	s.logger.Info().
		Str("transport", "grpc").
		Msg("stopping server")

	s.server.GracefulStop()
}

func newGrpcServer(logger zerolog.Logger, reg registration.Handler) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    5 * time.Minute,
				Timeout: 10 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             5 * time.Second,
				PermitWithoutStream: true,
			},
		),
		grpc.ChainUnaryInterceptor(
			makeLoggingInterceptor(logger),
			makeAuthInterceptor(reg),
		),
	)

	return grpcServer
}

func makeLoggingInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var correlationID string
		m, ok := metadata.FromIncomingContext(ctx)
		if ok {
			correlationID = getMetadataValue(m, "correlation_id")
		}

		correlationCtx := context.WithValue(ctx, CorrelationIDKey{}, correlationID)

		begin := time.Now()

		resp, err := handler(correlationCtx, req)

		logger := logger.With().
			Str("method", info.FullMethod).
			Str("correlation_id", correlationID).
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			logger.Error().Err(err).Msg("")
			return resp, err
		}

		logger.Info().Msg("called")

		return resp, err
	}
}

func makeAuthInterceptor(reg registration.Handler) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		m, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "identity metadata was not provided")
		}

		identity := application.Identity{
			VoterExternalId:        getMetadataValue(m, "voter_external_id"),
			VoterName:              getMetadataValue(m, "voter_name"),
			OrganizationExternalId: getMetadataValue(m, "org_external_id"),
			OrganizationName:       getMetadataValue(m, "org_name"),
		}

		voter, err := reg.RegisterIdentity(identity)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(context.WithValue(ctx, application.Voter{}, voter), req)
	}
}

func getMetadataValue(m metadata.MD, key string) string {
	value, ok := m[key]
	if !ok || len(value) < 1 {
		return ""
	}
	return value[0]
}

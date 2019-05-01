package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cheapRoc/grpc-zerolog"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/jukeizu/selection/api/protobuf-spec/selectionpb"
	"github.com/jukeizu/voting/api/protobuf-spec/votingpb"
	"github.com/jukeizu/voting/internal/startup"
	"github.com/jukeizu/voting/pkg/treediagram"
	"github.com/jukeizu/voting/pkg/voting"
	"github.com/jukeizu/voting/pkg/voting/ballot"
	"github.com/jukeizu/voting/pkg/voting/poll"
	"github.com/jukeizu/voting/pkg/voting/session"
	"github.com/jukeizu/voting/pkg/voting/voter"
	"github.com/oklog/run"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/shawntoffel/gossage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

var Version = ""

var (
	flagMigrate = false
	flagVersion = false
	flagServer  = false
	flagHandler = false

	grpcPort                = "50052"
	httpPort                = "10002"
	dbAddress               = "root@localhost:26257"
	serviceAddress          = "localhost:" + grpcPort
	selectionServiceAddress = "localhost:" + grpcPort
)

func parseConfig() {
	flag.StringVar(&grpcPort, "grpc.port", grpcPort, "grpc port for server")
	flag.StringVar(&httpPort, "http.port", httpPort, "http port for handler")
	flag.StringVar(&dbAddress, "db", dbAddress, "Database connection address")
	flag.StringVar(&serviceAddress, "service.addr", serviceAddress, "address of service if not local")
	flag.StringVar(&selectionServiceAddress, "selection.addr", selectionServiceAddress, "address of selection service if not local")
	flag.BoolVar(&flagServer, "server", false, "Run as server")
	flag.BoolVar(&flagHandler, "handler", false, "Run as handler")
	flag.BoolVar(&flagMigrate, "migrate", false, "Run db migrations")
	flag.BoolVar(&flagVersion, "v", false, "version")

	flag.Parse()
}

func main() {
	parseConfig()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("instance", xid.New().String()).
		Str("component", "voting").
		Str("version", Version).
		Logger()

	grpcLoggerV2 := grpczerolog.New(logger.With().Str("transport", "grpc").Logger())
	grpclog.SetLoggerV2(grpcLoggerV2)

	if !flagServer && !flagHandler {
		flagServer = true
		flagHandler = true
	}

	pollRepository, err := poll.NewRepository(dbAddress)
	if err != nil {
		logger.Error().Err(err).Caller().Msg("could not create poll repository")
		os.Exit(1)
	}

	sessionRepository, err := session.NewRepository(dbAddress)
	if err != nil {
		logger.Error().Err(err).Caller().Msg("could not create session repository")
		os.Exit(1)
	}

	voterRepository, err := voter.NewRepository(dbAddress)
	if err != nil {
		logger.Error().Err(err).Caller().Msg("could not create voter repository")
		os.Exit(1)
	}

	ballotRepository, err := ballot.NewRepository(dbAddress)
	if err != nil {
		logger.Error().Err(err).Caller().Msg("could not create ballot repository")
		os.Exit(1)
	}

	if flagMigrate {
		gossage.Logger = func(format string, a ...interface{}) {
			msg := fmt.Sprintf(format, a...)
			logger.Info().Str("component", "migrator").Msg(msg)
		}

		err = pollRepository.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not migrate poll repository")
			os.Exit(1)
		}

		err = sessionRepository.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not migrate session repository")
			os.Exit(1)
		}

		err = voterRepository.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not migrate voter repository")
			os.Exit(1)
		}

		err = ballotRepository.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not migrate ballot repository")
			os.Exit(1)
		}
	}

	clientConn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	if err != nil {
		logger.Error().Err(err).Str("serviceAddress", serviceAddress).Msg("could not dial service address")
		os.Exit(1)
	}

	selectionClientConn, err := grpc.Dial(selectionServiceAddress, grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	if err != nil {
		logger.Error().Err(err).Str("serviceAddress", serviceAddress).Msg("could not dial service address")
		os.Exit(1)
	}

	g := run.Group{}

	if flagServer {
		grpcServer := newGrpcServer(logger)

		server := startup.NewServer(logger, grpcServer)

		pollService := poll.NewDefaultService(logger, pollRepository)
		sessionService := session.NewDefaultService(logger, sessionRepository)
		voterService := voter.NewDefaultService(logger, voterRepository)
		ballotService := ballot.NewDefaultService(logger, ballotRepository)
		votingService := voting.NewDefaultService(logger, pollService, sessionService, voterService, ballotService)
		votingService = voting.NewValidationService(logger, votingService, pollService, voterService)

		votingServer := voting.NewGrpcServer(votingService)

		votingpb.RegisterVotingServer(grpcServer, votingServer)

		grpcAddr := ":" + grpcPort

		g.Add(func() error {
			return server.Start(grpcAddr)
		}, func(error) {
			server.Stop()
		})
	}

	if flagHandler {
		client := votingpb.NewVotingClient(clientConn)
		selectionClient := selectionpb.NewSelectionClient(selectionClientConn)

		httpAddr := ":" + httpPort

		handler := treediagram.NewHandler(logger, client, selectionClient, httpAddr)

		g.Add(func() error {
			return handler.Start()
		}, func(error) {
			err := handler.Stop()
			if err != nil {
				logger.Error().Err(err).Caller().Msg("couldn't stop handler")
			}
		})
	}

	cancel := make(chan struct{})
	g.Add(func() error {
		return interrupt(cancel)
	}, func(error) {
		close(cancel)
	})

	logger.Info().Err(g.Run()).Msg("stopped")
}

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-cancel:
		return errors.New("stopping")
	case sig := <-c:
		return fmt.Errorf("%s", sig)
	}
}

func newGrpcServer(logger zerolog.Logger) *grpc.Server {
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
		startup.LoggingInterceptor(logger),
	)

	return grpcServer
}

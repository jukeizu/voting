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
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/jukeizu/voting/poll"
	"github.com/jukeizu/voting/registration"
	registrationmediator "github.com/jukeizu/voting/registration/mediator"
	registrationrepository "github.com/jukeizu/voting/registration/repository"
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

	grpcPort       = "50052"
	httpPort       = "10002"
	dbAddress      = "root@localhost:26257"
	serviceAddress = "localhost:" + grpcPort
)

func parseConfig() {
	flag.StringVar(&grpcPort, "grpc.port", grpcPort, "grpc port for server")
	flag.StringVar(&httpPort, "http.port", httpPort, "http port for handler")
	flag.StringVar(&dbAddress, "db", dbAddress, "Database connection address")
	flag.StringVar(&serviceAddress, "service.addr", serviceAddress, "address of service if not local")
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

	if flagMigrate {
		gossage.Logger = func(format string, a ...interface{}) {
			msg := fmt.Sprintf(format, a...)
			logger.Info().Str("component", "migrator").Msg(msg)
		}

		pollRepository, err := poll.NewRepository(dbAddress)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not create poll repository")
			os.Exit(1)
		}

		err = pollRepository.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not migrate poll repository")
			os.Exit(1)
		}

		registrationRepository, err := registrationrepository.NewRepository(dbAddress)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not create registration repository")
			os.Exit(1)
		}

		err = registrationRepository.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not migrate registration repository")
			os.Exit(1)
		}
	}

	g := run.Group{}

	if flagServer {
		pollRepository, err := poll.NewRepository(dbAddress)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not create poll repository")
			os.Exit(1)
		}

		registrationRepository, err := registrationrepository.NewRepository(dbAddress)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not create registration repository")
			os.Exit(1)
		}

		grpcServer := newGrpcServer(logger)
		server := NewServer(logger, grpcServer)

		pollServer := poll.NewServer(logger, pollRepository)
		pollpb.RegisterPollsServer(grpcServer, pollServer)

		registrationMediator := registrationmediator.New()
		registrationServer := registration.NewServer(logger, registrationRepository)
		registrationpb.RegisterRegistrationServer(grpcServer, registrationServer)

		grpcAddr := ":" + grpcPort

		g.Add(func() error {
			return server.Start(grpcAddr)
		}, func(error) {
			server.Stop()
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
		LoggingInterceptor(logger),
	)

	return grpcServer
}

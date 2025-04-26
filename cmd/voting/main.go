package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	grpczerolog "github.com/cheapRoc/grpc-zerolog"
	"github.com/jukeizu/voting/internal/application/counting"
	"github.com/jukeizu/voting/internal/application/polling"
	"github.com/jukeizu/voting/internal/application/registration"
	"github.com/jukeizu/voting/internal/application/voting"
	"github.com/jukeizu/voting/internal/database"
	"github.com/jukeizu/voting/internal/infrastructure"
	"github.com/rs/zerolog"
	"github.com/shawntoffel/gossage"
	"google.golang.org/grpc/grpclog"
)

var Version = ""

const (
	connectionStringEnvVarName = "VOTING_DB_CONNECTION_STRING"
	defaultConnectionString    = "postgresql://postgres:password@localhost:5432/%s?sslmode=disable"
	consoleLogFormat           = "console"
)

var (
	flagMigrate = false
	flagVersion = false
	flagDebug   = false

	grpcPort      = "50052"
	flagLogFormat = consoleLogFormat
)

func init() {
	flag.StringVar(&grpcPort, "grpc.port", grpcPort, "grpc port for server")
	flag.BoolVar(&flagMigrate, "migrate", flagMigrate, "Run db migrations")
	flag.BoolVar(&flagDebug, "D", false, "enable debug logging")
	flag.StringVar(&flagLogFormat, "log-format", flagLogFormat, "log format: console, json")
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.Parse()
}

func main() {
	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if flagDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("component", "voting").
		Str("version", Version).
		Logger()

	if strings.EqualFold(flagLogFormat, consoleLogFormat) {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	grpcLoggerV2 := grpczerolog.New(logger.With().Str("transport", "grpc").Logger())
	grpclog.SetLoggerV2(grpcLoggerV2)

	gossage.Logger = func(format string, a ...interface{}) {
		msg := fmt.Sprintf(format, a...)
		logger.Info().Str("component", "migrator").Msg(msg)
	}

	dbAddress := readSecretEnvOrDefault(connectionStringEnvVarName, defaultConnectionString)
	db, err := database.New(dbAddress, flagMigrate)
	if err != nil {
		logger.Error().Err(err).Caller().Msg("could not open database")
		os.Exit(1)
	}

	regRepo := registration.NewRepository(db)
	regHandler := registration.NewValidatingHandler(registration.NewHandler(regRepo))

	pollRepo := polling.NewRepository(db)
	pollHandler := polling.NewValidatingHandler(polling.NewHandler(pollRepo), pollRepo)

	votingRepo := voting.NewRepository(db)
	votingHandler := voting.NewValidatingHandler(voting.NewHandler(votingRepo, pollRepo), votingRepo, pollRepo)

	countingRepo := counting.NewRepository(db)
	countingHandler := counting.NewValidatingHandler(counting.NewHandler(countingRepo, pollRepo), countingRepo, pollRepo)

	server := infrastructure.NewGrpcServer(
		logger,
		regHandler,
		pollHandler,
		votingHandler,
		countingHandler,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = server.Start(ctx, ":"+grpcPort)
	if err != nil {
		logger.Error().Err(err).Msg("failed to start server")
	}
}

func readSecretEnvOrDefault(name string, defaultValue string) string {
	env := os.Getenv(name)
	if len(env) > 0 {
		return env
	}

	file := os.Getenv(name + "_FILE")
	if len(file) < 1 {
		return defaultValue
	}

	bytes, err := os.ReadFile(file)
	if err != nil {
		return ""
	}

	return string(bytes)
}

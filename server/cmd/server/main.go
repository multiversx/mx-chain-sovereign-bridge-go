package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/core/closing"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/file"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/cmd/config"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var log = logger.GetOrCreate("sov-bridge-sender")

const (
	retrialTimeServe = 1
	envGRPCPort      = "GRPC_PORT"
	logsPath         = "logs"
	logsPrefix       = "sov-bridge-sender"
	logLifeSpanMb    = 1024   //# 1GB
	logLifeSpanSec   = 432000 // 5 days
)

func main() {
	app := cli.NewApp()
	app.Name = "Sovereign bridge tx server"
	app.Action = startServer
	app.Flags = []cli.Flag{
		logLevel,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startServer(ctx *cli.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	logFile, err := initializeLogger(ctx)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	bridgeServer, err := server.NewSovereignBridgeTxServer()
	if err != nil {
		return err
	}

	sovereign.RegisterBridgeTxSenderServer(grpcServer, bridgeServer)
	log.Info("starting server...")

	go func() {
		for {
			if err = grpcServer.Serve(listener); err != nil {
				log.LogIfError(err)
				time.Sleep(retrialTimeServe * time.Second)
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interrupt
	log.Info("closing app at user's signal")

	grpcServer.Stop()

	if !check.IfNilReflect(logFile) {
		err = logFile.Close()
		log.LogIfError(err)
	}

	return nil
}

func loadConfig() (*config.ServerConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	grpcPort := os.Getenv(envGRPCPort)

	log.Info("loaded config", "grpc port", grpcPort)

	return &config.ServerConfig{
		GRPCPort: grpcPort,
	}, nil
}

func initializeLogger(ctx *cli.Context) (closing.Closer, error) {
	logLevelFlagValue := ctx.GlobalString(logLevel.Name)
	err := logger.SetLogLevel(logLevelFlagValue)
	if err != nil {
		return nil, err
	}

	withLogFile := ctx.GlobalBool(logSaveFile.Name)
	if !withLogFile {
		return nil, nil
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.LogIfError(err)
		workingDir = ""
	}

	fileLogging, err := file.NewFileLogging(file.ArgsFileLogging{
		WorkingDir:      workingDir,
		DefaultLogsPath: logsPath,
		LogFilePrefix:   logsPrefix,
	})
	if err != nil {
		return nil, fmt.Errorf("%w creating a log file", err)
	}

	err = fileLogging.ChangeFileLifeSpan(
		time.Second*time.Duration(logLifeSpanSec),
		uint64(logLifeSpanMb),
	)
	if err != nil {
		return nil, err
	}

	disableAnsi := ctx.GlobalBool(disableAnsiColor.Name)
	err = removeANSIColorsForLoggerIfNeeded(disableAnsi)
	if err != nil {
		return nil, err
	}

	return fileLogging, nil
}

func removeANSIColorsForLoggerIfNeeded(disableAnsi bool) error {
	if !disableAnsi {
		return nil
	}

	err := logger.RemoveLogObserver(os.Stdout)
	if err != nil {
		return err
	}

	return logger.AddLogObserver(os.Stdout, &logger.PlainFormatter{})
}

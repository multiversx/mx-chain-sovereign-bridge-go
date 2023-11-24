package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/config"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var log = logger.GetOrCreate("server")

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

	err = initializeLogger(ctx)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	bridgeServer, err := server.NewServer()
	if err != nil {
		return err
	}

	sovereign.RegisterBridgeTxSenderServer(grpcServer, bridgeServer)

	log.Info("starting server...")

	go func() {
		for {
			if err = grpcServer.Serve(listener); err != nil {
				log.LogIfError(err)
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interrupt
	log.Info("closing app at user's signal")

	grpcServer.Stop()
	return nil
}

func loadConfig() (*config.ServerConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	grpcPort := os.Getenv("GRPC_PORT")

	log.Info("loaded config", "grpc port", grpcPort)

	return &config.ServerConfig{
		GRPCPort: grpcPort,
	}, nil
}

func initializeLogger(ctx *cli.Context) error {
	logLevelFlagValue := ctx.GlobalString(logLevel.Name)
	return logger.SetLogLevel(logLevelFlagValue)
}

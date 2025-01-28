package client

import (
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	logger "github.com/multiversx/mx-chain-logger-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/multiversx/mx-chain-sovereign-bridge-go/cert"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/client/config"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/client/disabled"
)

const (
	waitTime = 5
)

var log = logger.GetOrCreate("client")

// CreateClient creates a grpc client with retries
func CreateClient(cfg *config.ClientConfig) (ClientHandler, error) {
	if !cfg.Enabled {
		return disabled.NewClient(), nil
	}

	dialTarget := fmt.Sprintf("%s:%s", cfg.GRPCHost, cfg.GRPCPort)
	conn, err := connectWithRetries(dialTarget, cfg.CertificateCfg)
	if err != nil {
		return nil, err
	}

	bridgeClient := sovereign.NewBridgeTxSenderClient(conn)
	return NewClient(bridgeClient, conn)
}

func connectWithRetries(host string, cfg cert.FileCfg) (GRPCConn, error) {
	tlsConfig, err := cert.LoadTLSClientConfig(cfg)
	if err != nil {
		return nil, err
	}

	for i := 0; ; i++ {
		tlsCredentials := credentials.NewTLS(tlsConfig)
		cc, errConnection := grpc.Dial(host, grpc.WithTransportCredentials(tlsCredentials))
		if errConnection == nil {
			return cc, errConnection
		}

		time.Sleep(time.Second * waitTime)

		log.Warn("could not establish connection, retrying",
			"error", errConnection,
			"host", host,
			"retries", i+1)
	}
}

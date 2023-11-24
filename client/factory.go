package client

import (
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/client/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	maxConnectionRetrials = 100
	waitTime              = 5
)

var log = logger.GetOrCreate("client")

// CreateClient creates a grpc client with retrials
func CreateClient(cfg *config.ClientConfig) (ClientHandler, error) {
	dialTarget := fmt.Sprintf("%s:%s", cfg.GRPCHost, cfg.GRPCPort)
	conn, err := connectWithRetrials(dialTarget)
	if err != nil {
		return nil, err
	}

	bridgeClient := sovereign.NewBridgeTxSenderClient(conn)
	return NewClient(bridgeClient, conn)
}

func connectWithRetrials(host string) (GRPCConn, error) {
	credentials := insecure.NewCredentials()
	opts := grpc.WithTransportCredentials(credentials)

	for i := 0; i < maxConnectionRetrials; i++ {
		cc, err := grpc.Dial(host, opts)
		if err == nil {
			return cc, err
		}

		time.Sleep(time.Second * waitTime)

		log.Warn("could not establish connection, retrying",
			"error", err,
			"host", host,
			"retrial", i+1)
	}

	return nil, errCannotOpenConnection
}

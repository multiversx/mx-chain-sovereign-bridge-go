package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/client"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/client/config"
	"github.com/urfave/cli"
)

var log = logger.GetOrCreate("server")

func main() {
	app := cli.NewApp()
	app.Name = "GRPC Client"
	app.Action = startClient
	app.Flags = []cli.Flag{
		logLevel,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startClient(ctx *cli.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	err = initializeLogger(ctx)
	if err != nil {
		return err
	}

	log.Info("starting client...")

	bridgeClient, err := client.CreateClient(cfg)
	if err != nil {
		return err
	}

	defer func() {
		err = bridgeClient.Close()
		log.LogIfError(err)
	}()

	return sendData(bridgeClient)
}

func sendData(bridgeClient client.ClientHandler) error {

	txHashes := make([]string, 0)
	mut := sync.RWMutex{}

	numBridgeOps := 5
	for i := 0; i < numBridgeOps; i++ {
		hash := []byte(fmt.Sprintf("hash_%d", i))
		log.Info("sending data", "hash", hash)

		go func() {
			res, errSend := bridgeClient.Send(context.Background(), &sovereign.BridgeOperations{
				Data: []*sovereign.BridgeOutGoingData{
					{
						Hash: hash,
						OutGoingOperations: map[string][]byte{
							"fc07": []byte("bridgeOp"),
						},
					},
				},
			})
			if errSend != nil {
				log.Error("error sending bridge data", "error", errSend)
				return
			}

			logTxHashes(res.TxHashes)

			mut.Lock()
			txHashes = append(txHashes, res.TxHashes...)
			mut.Unlock()
		}()
	}

	time.Sleep(time.Second * 50)

	numSentTxs := len(txHashes)
	expectedNumBridgeTxs := 2 * numBridgeOps
	if numSentTxs != expectedNumBridgeTxs {
		return fmt.Errorf("did not send all txs; expected num send txs: %d, received: %d",
			expectedNumBridgeTxs, numSentTxs)
	}

	return nil
}

func logTxHashes(txHashes []string) {
	for _, txHash := range txHashes {
		log.Info("received", "tx hash", txHash)
	}
}

func loadConfig() (*config.ClientConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	grpcHost := os.Getenv("GRPC_HOST")
	grpcPort := os.Getenv("GRPC_PORT")

	log.Info("loaded config", "grpc host", grpcHost)
	log.Info("loaded config", "grpc port", grpcPort)

	return &config.ClientConfig{
		GRPCHost: grpcHost,
		GRPCPort: grpcPort,
	}, nil
}

func initializeLogger(ctx *cli.Context) error {
	logLevelFlagValue := ctx.GlobalString(logLevel.Name)
	return logger.SetLogLevel(logLevelFlagValue)
}

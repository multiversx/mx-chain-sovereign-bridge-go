package server

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/data/sovereign/dto"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("server")

type server struct {
	txSenders map[dto.ChainID]TxSender
	*sovereign.UnimplementedBridgeTxSenderServer
}

// NewSovereignBridgeTxServer creates a new sovereign bridge operations server. This server receives bridge data operations from
// sovereign nodes and sends transactions to main chain.
func NewSovereignBridgeTxServer(txSenders map[dto.ChainID]TxSender) (*server, error) {
	if len(txSenders) == 0 {
		return nil, errNilTxSender
	}

	for chainID, txSender := range txSenders {
		if check.IfNil(txSender) {
			return nil, fmt.Errorf("%w for chain id: %s", errNilTxSender, chainID.String())
		}
	}

	return &server{
		txSenders: txSenders,
	}, nil
}

// Send should handle receiving data bridge operations from sovereign shard and forward transactions to main chain
func (s *server) Send(ctx context.Context, data *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
	dataToSendPerChain := getDataToSendPerChain(data)

	allHashes := make([]string, 0)
	for chainID, dataToSend := range dataToSendPerChain {
		txSender, found := s.txSenders[chainID]
		if !found {
			log.Error("received data to bridge for unknown chain id", "chainID", chainID.String())
			continue
		}

		hashes, err := txSender.SendTxs(ctx, dataToSend)
		if err != nil {
			return nil, err
		}

		logTxHashes(chainID, hashes)

		allHashes = append(allHashes, hashes...)
	}

	slices.SortStableFunc(allHashes, func(a, b string) int {
		return strings.Compare(a, b)
	})

	return &sovereign.BridgeOperationsResponse{
		TxHashes: allHashes,
	}, nil
}

func getDataToSendPerChain(data *sovereign.BridgeOperations) map[dto.ChainID]*sovereign.BridgeOperations {
	dataToSendPerChain := make(map[dto.ChainID]*sovereign.BridgeOperations)
	for _, dta := range data.Data {
		chainID := dto.ChainID(dta.ChainID)

		if ops, ok := dataToSendPerChain[chainID]; ok {
			ops.Data = append(ops.Data, dta)
		} else {
			dataToSendPerChain[chainID] = &sovereign.BridgeOperations{
				Data: []*sovereign.BridgeOutGoingData{dta},
			}
		}
	}

	return dataToSendPerChain
}

func logTxHashes(chainID dto.ChainID, hashes []string) {
	for _, hash := range hashes {
		log.Info("sent tx", "chain", chainID.String(), "hash", hash)
	}
}

// IsInterfaceNil checks if the underlying pointer is nil
func (s *server) IsInterfaceNil() bool {
	return s == nil
}

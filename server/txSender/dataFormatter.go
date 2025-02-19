package txSender

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/hashing"
)

const (
	registerBridgeOpsPrefix = "registerBridgeOps"
	executeBridgeOpsPrefix  = "executeBridgeOps"
)

type dataFormatter struct {
	dataFormatterHandlers map[int32]txDataFormatter
}

// NewDataFormatter creates a sovereign bridge tx data formatter
func NewDataFormatter(hasher hashing.Hasher) (*dataFormatter, error) {
	if check.IfNil(hasher) {
		return nil, core.ErrNilHasher
	}

	return &dataFormatter{
		dataFormatterHandlers: map[int32]txDataFormatter{
			int32(block.OutGoingMbTx):                 &dataFormatterDepositTokens{hasher: hasher},
			int32(block.OutGoingMbChangeValidatorSet): &dataFormatterValidatorSetChange{},
		},
	}, nil
}

// CreateTxsData creates txs data for bridge operations
func (df *dataFormatter) CreateTxsData(data *sovereign.BridgeOperations) [][]byte {
	txsData := make([][]byte, 0)
	if data == nil {
		return txsData
	}

	for _, bridgeData := range data.Data {
		log.Debug("creating tx data",
			"type", block.OutGoingMBType(bridgeData.Type).String(),
			"bridge op hash", bridgeData.Hash,
			"no. of operations", len(bridgeData.OutGoingOperations),
		)

		handler, found := df.dataFormatterHandlers[bridgeData.Type]
		if !found {
			log.Error("received unknown bridge data", "type", bridgeData.Type)
			continue
		}

		newTxsData, err := handler.createTxsData(bridgeData)
		if err != nil {
			log.Error("could not create txs data",
				"error", err,
				"hash", bridgeData.Hash,
				"type", block.OutGoingMBType(bridgeData.Type).String())
			continue
		}

		txsData = append(txsData, newTxsData...)
	}

	return txsData
}

// IsInterfaceNil checks if the underlying pointer is nil
func (df *dataFormatter) IsInterfaceNil() bool {
	return df == nil
}

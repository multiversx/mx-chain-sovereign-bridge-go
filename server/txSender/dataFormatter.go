package txSender

import (
	"bytes"
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/hashing"
)

const (
	registerBridgeOpsPrefix = "registerBridgeOps"
	executeBridgeOpsPrefix  = "executeBridgeOps"
)

type dataFormatter struct {
	hasher hashing.Hasher
}

// NewDataFormatter creates a sovereign bridge tx data formatter
func NewDataFormatter() *dataFormatter {
	return &dataFormatter{}
}

// CreateTxsData creates txs data for bridge operations
func (df *dataFormatter) CreateTxsData(data *sovereign.BridgeOperations) [][]byte {
	txsData := make([][]byte, 0)
	if data == nil {
		return txsData
	}

	for _, bridgeData := range data.Data {
		log.Debug("creating tx data", "bridge op hash", bridgeData.Hash)

		registerBridgeOpData := df.createRegisterBridgeOperationsData(bridgeData)
		if len(registerBridgeOpData) != 0 {
			txsData = append(txsData, df.createRegisterBridgeOperationsData(bridgeData))
		}

		txsData = append(txsData, createBridgeOperationsData(bridgeData.Hash, bridgeData.OutGoingOperations)...)
	}

	return txsData
}

func (df *dataFormatter) createRegisterBridgeOperationsData(bridgeData *sovereign.BridgeOutGoingData) []byte {
	args := []byte(registerBridgeOpsPrefix +
		"@" + hex.EncodeToString(bridgeData.AggregatedSignature) +
		"@" + hex.EncodeToString(bridgeData.Hash))
	for _, operation := range bridgeData.OutGoingOperations {
		args = append(args, "@"+hex.EncodeToString(operation.Hash)...)
	}

	// unconfirmed operation, should not register it, only resend it
	computedHashOfHashes := df.hasher.Compute(string(hashes))
	if !bytes.Equal(hashOfHashes, computedHashOfHashes) {
		return nil
	}

	return args
}

func createBridgeOperationsData(hashOfHashes []byte, outGoingOperations []*sovereign.OutGoingOperation) [][]byte {
	ret := make([][]byte, 0)
	for _, operation := range outGoingOperations {
		currBridgeOp := []byte(executeBridgeOpsPrefix + "@" + hex.EncodeToString(hashOfHashes) + "@")
		currBridgeOp = append(currBridgeOp, hex.EncodeToString(operation.Data)...)

		ret = append(ret, currBridgeOp)
	}

	return ret
}

// IsInterfaceNil checks if the underlying pointer is nil
func (df *dataFormatter) IsInterfaceNil() bool {
	return df == nil
}

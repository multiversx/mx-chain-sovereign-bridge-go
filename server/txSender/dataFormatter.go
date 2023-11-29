package txSender

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

const (
	registerBridgeOpsPrefix = "registerBridgeOps"
	executeBridgeOpPrefix   = "executeBridgeOp"
)

type dataFormatter struct {
}

// NewDataFormatter creates a sovereign bridge tx data formatter
func NewDataFormatter() *dataFormatter {
	return &dataFormatter{}
}

// CreateTxsData creates txs data for bridge operations
func (df *dataFormatter) CreateTxsData(data *sovereign.BridgeOperations) [][]byte {
	txsData := make([][]byte, 0)

	for _, bridgeData := range data.Data {
		log.Debug("creating tx data", "bridge op hash", bridgeData.Hash)
		txsData = append(txsData, createRegisterBridgeOperationsData(bridgeData))
		txsData = append(txsData, createBridgeOperationsData(bridgeData.OutGoingOperations)...)
	}

	return txsData
}

func createRegisterBridgeOperationsData(bridgeData *sovereign.BridgeOutGoingData) []byte {
	registerBridgeOpTxData := []byte(
		registerBridgeOpsPrefix + "@" +
			hex.EncodeToString(bridgeData.LeaderSignature) + "@" +
			hex.EncodeToString(bridgeData.AggregatedSignature) + "@" +
			hex.EncodeToString(bridgeData.Hash))

	listOfOps := make([]byte, 0, len(bridgeData.OutGoingOperations))
	for operationHash := range bridgeData.OutGoingOperations {
		listOfOps = append(listOfOps, []byte("@")...)
		listOfOps = append(listOfOps, []byte(operationHash)...)
	}

	return append(registerBridgeOpTxData, listOfOps...)
}

func createBridgeOperationsData(outGoingOperations map[string][]byte) [][]byte {
	ret := make([][]byte, 0)
	for operationHash, bridgeOpData := range outGoingOperations {
		currBridgeOp := []byte(executeBridgeOpPrefix + "@" + operationHash + "@")
		currBridgeOp = append(currBridgeOp, bridgeOpData...)

		ret = append(ret, currBridgeOp)
	}

	return ret
}

// IsInterfaceNil checks if the underlying pointer is nil
func (df *dataFormatter) IsInterfaceNil() bool {
	return df == nil
}

package txSender

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

const (
	registerBridgeOpsPrefix = "registerBridgeOps"
	executeBridgeOpsPrefix  = "executeBridgeOps"
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
	if data == nil {
		return txsData
	}

	for _, bridgeData := range data.Data {
		log.Debug("creating tx data", "bridge op hash", bridgeData.Hash)
		txsData = append(txsData, createRegisterBridgeOperationsData(bridgeData))
		txsData = append(txsData, createBridgeOperationsData(bridgeData.Hash, bridgeData.OutGoingOperations)...)
	}

	return txsData
}

func createRegisterBridgeOperationsData(bridgeData *sovereign.BridgeOutGoingData) []byte {
	registerBridgeOpTxData := []byte(
		registerBridgeOpsPrefix +
			"@" + hex.EncodeToString(bridgeData.AggregatedSignature) +
			"@" + hex.EncodeToString(bridgeData.Hash))

	for _, operation := range bridgeData.OutGoingOperations {
		registerBridgeOpTxData = append(registerBridgeOpTxData, "@"+hex.EncodeToString(operation.Hash)...)
	}

	return registerBridgeOpTxData
}

func createBridgeOperationsData(hashOfHashes []byte, outGoingOperations []*sovereign.OutGoingOperation) [][]byte {
	executeBridgeOpsTxData := make([][]byte, 0)
	for _, operation := range outGoingOperations {
		bridgeOpTxData := []byte(
			executeBridgeOpsPrefix +
				"@" + hex.EncodeToString(hashOfHashes) +
				"@" + hex.EncodeToString(operation.Data))

		executeBridgeOpsTxData = append(executeBridgeOpsTxData, bridgeOpTxData)
	}

	return executeBridgeOpsTxData
}

// IsInterfaceNil checks if the underlying pointer is nil
func (df *dataFormatter) IsInterfaceNil() bool {
	return df == nil
}

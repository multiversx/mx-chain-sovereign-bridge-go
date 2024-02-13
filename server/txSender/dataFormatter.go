package txSender

import (
	"encoding/hex"
	logger "github.com/multiversx/mx-chain-logger-go"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

const (
	registerBridgeOpsPrefix = "registerBridgeOps"
	executeBridgeOpPrefix   = "executeBridgeOp"
)

var (
	log = logger.GetOrCreate("mx-chain-sovereign-bridge-go")
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
	for _, operation := range bridgeData.OutGoingOperations {
		listOfOps = append(listOfOps, []byte("@")...)
		listOfOps = append(listOfOps, []byte(hex.EncodeToString(operation.Hash))...)
	}

	return append(registerBridgeOpTxData, listOfOps...)
}

func createBridgeOperationsData(outGoingOperations []*sovereign.OutGoingOperation) [][]byte {
	ret := make([][]byte, 0)
	for _, operation := range outGoingOperations {
		currBridgeOp := []byte(executeBridgeOpPrefix + "@" + hex.EncodeToString(operation.Hash) + "@")
		currBridgeOp = append(currBridgeOp, operation.Data...)

		ret = append(ret, currBridgeOp)
	}

	return ret
}

// IsInterfaceNil checks if the underlying pointer is nil
func (df *dataFormatter) IsInterfaceNil() bool {
	return df == nil
}

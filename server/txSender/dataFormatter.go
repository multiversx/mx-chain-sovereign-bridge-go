package txSender

import (
	"encoding/hex"
	"fmt"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"strconv"
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
		txsData = append(txsData, createBridgeOperationsData(bridgeData.Hash, bridgeData.OutGoingOperations))
	}

	return txsData
}

func createRegisterBridgeOperationsData(bridgeData *sovereign.BridgeOutGoingData) []byte {
	hashOfHashes := bridgeData.Hash
	hashes := make([]byte, 0)
	for _, operation := range bridgeData.OutGoingOperations {
		hashes = append(hashes, valueToHexString(len(operation.Hash), 4)...)
		hashes = append(hashes, operation.Hash...)
	}

	return []byte(registerBridgeOpsPrefix + "@" +
		hex.EncodeToString(hashOfHashes) + "@" +
		hex.EncodeToString(hashes) + "@" +
		hex.EncodeToString(bridgeData.AggregatedSignature))
}

func createBridgeOperationsData(hashOfHashes []byte, outGoingOperations []*sovereign.OutGoingOperation) []byte {
	bridgeOps := []byte(executeBridgeOpsPrefix + "@")
	bridgeOps = append(bridgeOps, hex.EncodeToString(hashOfHashes)...)
	for _, operation := range outGoingOperations {
		bridgeOps = append(bridgeOps, "@"+hex.EncodeToString(operation.Data)...)
	}

	return bridgeOps
}

func valueToHexString(value int, size int) []byte {
	hexString := strconv.FormatInt(int64(value), 16)
	paddedHexString := fmt.Sprintf("%016s", hexString)

	decoded, _ := hex.DecodeString(paddedHexString[(size * 2):])
	return decoded
}

// IsInterfaceNil checks if the underlying pointer is nil
func (df *dataFormatter) IsInterfaceNil() bool {
	return df == nil
}

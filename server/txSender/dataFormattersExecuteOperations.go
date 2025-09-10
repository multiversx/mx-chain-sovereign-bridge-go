package txSender

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

func createExecuteDepositTokensBridgeOperationsData(hashOfHashes []byte, outGoingOperations []*sovereign.OutGoingOperation) [][]byte {
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

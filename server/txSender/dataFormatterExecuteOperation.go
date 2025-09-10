package txSender

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/hashing"
)

const registerBridgeOpsPrefix = "registerBridgeOps"

const (
	executeDepositBridgeOpsPrefix    = "executeBridgeOps"
	executeRegisterTokenPrefix       = "registerToken"
	executeRegisterValidatorPrefix   = "registerValidator"
	executeUnRegisterValidatorPrefix = "unRegisterValidator"
)

type dataFormatterExecuteOperation struct {
	hasher          hashing.Hasher
	executeOpPrefix string
}

func (df *dataFormatterExecuteOperation) createTxsData(bridgeData *sovereign.BridgeOutGoingData) ([][]byte, error) {
	txsData := make([][]byte, 0)
	registerBridgeOpData := df.createRegisterBridgeOperationsData(bridgeData)
	if len(registerBridgeOpData) != 0 {
		txsData = append(txsData, registerBridgeOpData)
	}

	return append(txsData, df.createExecuteDepositTokensBridgeOperationsData(bridgeData.Hash, bridgeData.OutGoingOperations)...), nil
}

func (df *dataFormatterExecuteOperation) createRegisterBridgeOperationsData(bridgeData *sovereign.BridgeOutGoingData) []byte {
	hashes := make([]byte, 0)
	hashesHexEncodedArgs := make([]byte, 0)
	for _, operation := range bridgeData.OutGoingOperations {
		hashesHexEncodedArgs = append(hashesHexEncodedArgs, "@"+hex.EncodeToString(operation.Hash)...)
		hashes = append(hashes, operation.Hash...)
	}

	// unconfirmed operation, should not register it, only resend it
	computedHashOfHashes := df.hasher.Compute(string(hashes))
	if !bytes.Equal(bridgeData.Hash, computedHashOfHashes) {
		return nil
	}

	registerBridgeOpData := []byte(registerBridgeOpsPrefix +
		"@" + hex.EncodeToString(bridgeData.AggregatedSignature) +
		"@" + hex.EncodeToString(bridgeData.Hash) +
		"@" + hex.EncodeToString(bridgeData.PubKeysBitmap) +
		"@" + hex.EncodeToString(uint32ToBytes(bridgeData.Epoch)))

	return append(registerBridgeOpData, hashesHexEncodedArgs...)
}

func uint32ToBytes(value uint32) []byte {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, value)
	return buff
}

func (df *dataFormatterExecuteOperation) createExecuteDepositTokensBridgeOperationsData(hashOfHashes []byte, outGoingOperations []*sovereign.OutGoingOperation) [][]byte {
	executeBridgeOpsTxData := make([][]byte, 0)
	for _, operation := range outGoingOperations {
		bridgeOpTxData := []byte(
			df.executeOpPrefix +
				"@" + hex.EncodeToString(hashOfHashes) +
				"@" + hex.EncodeToString(operation.Data))

		executeBridgeOpsTxData = append(executeBridgeOpsTxData, bridgeOpTxData)
	}

	return executeBridgeOpsTxData
}

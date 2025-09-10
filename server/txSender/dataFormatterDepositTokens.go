package txSender

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/hashing"
)

const (
	registerBridgeOpsPrefix = "registerBridgeOps"
	executeBridgeOpsPrefix  = "executeBridgeOps"
)

type createExecuteBridgeOperationsDataHandler func(hashOfHashes []byte, outGoingOperations []*sovereign.OutGoingOperation) [][]byte

type dataFormatterExecuteOperation struct {
	hasher                                hashing.Hasher
	createExecuteBridgeOperationsDataFunc createExecuteBridgeOperationsDataHandler
}

func (df *dataFormatterExecuteOperation) createTxsData(bridgeData *sovereign.BridgeOutGoingData) ([][]byte, error) {
	txsData := make([][]byte, 0)
	registerBridgeOpData := df.createRegisterBridgeOperationsData(bridgeData)
	if len(registerBridgeOpData) != 0 {
		txsData = append(txsData, registerBridgeOpData)
	}

	return append(txsData, df.createExecuteBridgeOperationsDataFunc(bridgeData.Hash, bridgeData.OutGoingOperations)...), nil
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

package txSender

import (
	"bytes"
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/data/typeConverters"
	"github.com/multiversx/mx-chain-core-go/data/typeConverters/uint64ByteSlice"
	"github.com/multiversx/mx-chain-core-go/hashing"
)

const (
	registerBridgeOpsPrefix = "registerBridgeOps"
	executeBridgeOpsPrefix  = "executeBridgeOps"
)

type dataFormatter struct {
	hasher        hashing.Hasher
	uIntConverter typeConverters.Uint64ByteSliceConverter
}

// NewDataFormatter creates a sovereign bridge tx data formatter
func NewDataFormatter(hasher hashing.Hasher) (*dataFormatter, error) {
	if check.IfNil(hasher) {
		return nil, core.ErrNilHasher
	}

	return &dataFormatter{
		hasher:        hasher,
		uIntConverter: uint64ByteSlice.NewBigEndianConverter(),
	}, nil
}

// CreateTxsData creates txs data for bridge operations
func (df *dataFormatter) CreateTxsData(data *sovereign.BridgeOperations) [][]byte {
	txsData := make([][]byte, 0)
	if data == nil {
		return txsData
	}

	for _, bridgeData := range data.Data {
		log.Debug("creating tx data", "bridge op hash", bridgeData.Hash, "no. of operations", len(bridgeData.OutGoingOperations))

		registerBridgeOpData := df.createRegisterBridgeOperationsData(bridgeData)
		if len(registerBridgeOpData) != 0 {
			txsData = append(txsData, registerBridgeOpData)
		}

		txsData = append(txsData, createBridgeOperationsData(bridgeData.Hash, bridgeData.OutGoingOperations)...)
	}

	return txsData
}

func (df *dataFormatter) createRegisterBridgeOperationsData(bridgeData *sovereign.BridgeOutGoingData) []byte {
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

	// TODO: Here, create a new uint32 converter or just use binary.Uin32 std library

	registerBridgeOpData := []byte(registerBridgeOpsPrefix +
		"@" + hex.EncodeToString(bridgeData.AggregatedSignature) +
		"@" + hex.EncodeToString(bridgeData.Hash))
	registerBridgeOpData = append(registerBridgeOpData, hashesHexEncodedArgs...)

	bridgeDataArgs := []byte(
		"@" + hex.EncodeToString(bridgeData.PubKeysBitmap) +
			"@" + hex.EncodeToString(df.uIntConverter.ToByteSlice(uint64(bridgeData.Epoch))))

	return append(registerBridgeOpData, bridgeDataArgs...)
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

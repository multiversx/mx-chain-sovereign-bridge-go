package txSender

import (
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

const changeValidatorSetPrefix = "changeValidatorSet"

type dataFormatterValidatorSetChange struct {
}

func newDataFormatterValidatorSetChange() *dataFormatterValidatorSetChange {
	return &dataFormatterValidatorSetChange{}
}

// createTxsData will format the data to the following format:
//
// changeValidatorSet@HashOfHashes@HashOfOperation@AggregatedBLSMultiSig@PubKeysBitMap@Epoch@list<allKeyIDsInNewEpoch>
func (df *dataFormatterValidatorSetChange) createTxsData(bridgeData *sovereign.BridgeOutGoingData) ([][]byte, error) {
	numOutGoingOperations := len(bridgeData.OutGoingOperations)
	if numOutGoingOperations != 1 {
		return nil, fmt.Errorf("%w, expected 1, got %d", errInvalidBridgeDataSetValidatorChange, numOutGoingOperations)
	}

	pubKeys, err := formatPubKeys(bridgeData.OutGoingOperations[0].Data)
	if err != nil {
		return nil, err
	}

	txData := []byte(changeValidatorSetPrefix +
		"@" + hex.EncodeToString(bridgeData.Hash) +
		"@" + hex.EncodeToString(bridgeData.OutGoingOperations[0].Hash) +
		"@" + hex.EncodeToString(bridgeData.AggregatedSignature) +
		"@" + hex.EncodeToString(bridgeData.PubKeysBitmap) +
		"@" + hex.EncodeToString(uint32ToBytes(bridgeData.Epoch)))

	return [][]byte{append(txData, pubKeys...)}, nil
}

func formatPubKeys(data []byte) ([]byte, error) {
	pubKeysBridgeData := &sovereign.BridgeOutGoingDataValidatorSetChange{
		PubKeyIDs: make([][]byte, 0),
	}
	err := proto.Unmarshal(data, pubKeysBridgeData)
	if err != nil {
		return nil, err
	}

	pubKeysHex := make([]byte, 0, len(pubKeysBridgeData.PubKeyIDs))
	for _, pubKey := range pubKeysBridgeData.PubKeyIDs {
		pubKeysHex = append(pubKeysHex, "@"+hex.EncodeToString(pubKey)...)
	}

	return pubKeysHex, nil
}

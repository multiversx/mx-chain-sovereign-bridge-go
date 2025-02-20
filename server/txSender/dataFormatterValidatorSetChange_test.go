package txSender

import (
	"encoding/hex"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/require"
)

func TestDataFormatterValidatorSetChange_createTxsData(t *testing.T) {
	dataFormatterValidators := newDataFormatterValidatorSetChange()

	pubKey1 := []byte("pk1")
	pubKey2 := []byte("pk2")
	validatorPubKeysData := &sovereign.BridgeOutGoingDataValidatorSetChange{
		PubKeyIDs: [][]byte{pubKey1, pubKey2},
	}
	validatorPubKeysDataBytes, err := proto.Marshal(validatorPubKeysData)
	require.Nil(t, err)

	bridgeData := &sovereign.BridgeOutGoingData{
		Type: int32(block.OutGoingMbChangeValidatorSet),
		Hash: []byte("hashOfHashes"),
		OutGoingOperations: []*sovereign.OutGoingOperation{
			{
				Hash: []byte("operationHash"),
				Data: validatorPubKeysDataBytes,
			},
		},
		AggregatedSignature: []byte("aggregatedSig"),
		PubKeysBitmap:       []byte("pubKeysBitmap"),
		Epoch:               4,
	}

	expectedTxData := []byte(changeValidatorSetPrefix +
		"@" + hex.EncodeToString([]byte("hashOfHashes")) +
		"@" + hex.EncodeToString([]byte("operationHash")) +
		"@" + hex.EncodeToString([]byte("aggregatedSig")) +
		"@" + hex.EncodeToString([]byte("pubKeysBitmap")) +
		"@" + "00000004" +
		"@" + hex.EncodeToString(pubKey1) +
		"@" + hex.EncodeToString(pubKey2))

	txData, err := dataFormatterValidators.createTxsData(bridgeData)
	require.Nil(t, err)
	require.Equal(t, [][]byte{expectedTxData}, txData)
}

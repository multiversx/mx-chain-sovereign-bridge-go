package txSender

import (
	"encoding/hex"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/stretchr/testify/require"
)

func TestNewDataFormatter(t *testing.T) {
	t.Parallel()

	df := NewDataFormatter()
	require.False(t, df.IsInterfaceNil())
}

func TestDataFormatter_CreateTxsData(t *testing.T) {
	t.Parallel()

	df := NewDataFormatter()

	t.Run("nil input, should return empty result", func(t *testing.T) {
		require.Empty(t, df.CreateTxsData(nil))
	})

	t.Run("empty data, should return empty result", func(t *testing.T) {
		require.Empty(t, df.CreateTxsData(&sovereign.BridgeOperations{Data: nil}))
	})

	t.Run("should work", func(t *testing.T) {
		bridgeDataHash1 := []byte("bridgeDataHash1")
		bridgeDataHash2 := []byte("bridgeDataHash2")

		aggregatedSig1 := []byte("aggregatedSig1")
		aggregatedSig2 := []byte("aggregatedSig2")

		leaderSig1 := []byte("leaderSig1")
		leaderSig2 := []byte("leaderSig2")

		opHash1 := []byte("outGoingOpHash1")
		opHash2 := []byte("outGoingOpHash2")
		opHash3 := []byte("outGoingOpHash3")

		bridgeOps := &sovereign.BridgeOperations{
			Data: []*sovereign.BridgeOutGoingData{
				{
					Hash: bridgeDataHash1,
					OutGoingOperations: []*sovereign.OutGoingOperation{
						{
							Hash: opHash1,
							Data: []byte("bridgeOp1"),
						},
						{
							Hash: opHash2,
							Data: []byte("bridgeOp2@bridgeOp22"),
						},
					},
					AggregatedSignature: aggregatedSig1,
					LeaderSignature:     leaderSig1,
				},
				{
					Hash: bridgeDataHash2,
					OutGoingOperations: []*sovereign.OutGoingOperation{
						{
							Hash: opHash3,
							Data: []byte("bridgeOp3"),
						},
					},
					AggregatedSignature: aggregatedSig2,
					LeaderSignature:     leaderSig2,
				},
			},
		}

		registerOp1 := []byte(
			registerBridgeOpsPrefix + "@" +
				hex.EncodeToString(leaderSig1) + "@" +
				hex.EncodeToString(aggregatedSig1) + "@" +
				hex.EncodeToString(bridgeDataHash1) + "@" +
				hex.EncodeToString(opHash1) + "@" +
				hex.EncodeToString(opHash2))
		execOp1 := []byte(executeBridgeOpPrefix + "@" +
			hex.EncodeToString(opHash1) + "@" +
			"bridgeOp1")
		execOp2 := []byte(executeBridgeOpPrefix + "@" +
			hex.EncodeToString(opHash2) + "@" +
			"bridgeOp2@bridgeOp22")

		registerOp2 := []byte(
			registerBridgeOpsPrefix + "@" +
				hex.EncodeToString(leaderSig2) + "@" +
				hex.EncodeToString(aggregatedSig2) + "@" +
				hex.EncodeToString(bridgeDataHash2) + "@" +
				hex.EncodeToString(opHash3))
		execOp3 := []byte(executeBridgeOpPrefix + "@" +
			hex.EncodeToString(opHash3) + "@" +
			"bridgeOp3")

		expectedTxsData := [][]byte{
			registerOp1,
			execOp1,
			execOp2,
			registerOp2,
			execOp3,
		}

		txsData := df.CreateTxsData(bridgeOps)
		require.Equal(t, expectedTxsData, txsData)
	})
}

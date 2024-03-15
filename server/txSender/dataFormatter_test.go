package txSender

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewDataFormatter(t *testing.T) {
	t.Parallel()

	t.Run("nil hasher, should fail", func(t *testing.T) {
		df, err := NewDataFormatter(nil)
		require.Equal(t, core.ErrNilHasher, err)
		require.Nil(t, df)
	})

	t.Run("should work", func(t *testing.T) {
		df, err := NewDataFormatter(&testscommon.HasherMock{})
		require.Nil(t, err)
		require.False(t, df.IsInterfaceNil())
	})
}

func TestDataFormatter_CreateTxsData(t *testing.T) {
	t.Parallel()

	t.Run("nil input, should return empty result", func(t *testing.T) {
		df, _ := NewDataFormatter(&testscommon.HasherMock{})
		require.Empty(t, df.CreateTxsData(nil))
	})

	t.Run("empty data, should return empty result", func(t *testing.T) {
		df, _ := NewDataFormatter(&testscommon.HasherMock{})
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
							Data: []byte("bridgeOp2"),
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

		computeHashCt := 0
		hasher := &testscommon.HasherMock{
			ComputeCalled: func(s string) []byte {
				defer func() {
					computeHashCt++
				}()

				switch computeHashCt {
				case 0:
					require.Equal(t, string(append(opHash1, opHash2...)), s)
					return bridgeDataHash1
				case 1:
					require.Equal(t, string(opHash3), s)
					return bridgeDataHash2
				default:
					require.Fail(t, "should have not compute another hash")
				}

				return nil
			},
		}
		df, _ := NewDataFormatter(hasher)

		registerOp1 := []byte(
			registerBridgeOpsPrefix + "@" +
				hex.EncodeToString(bridgeDataHash1) + "@" +
				fmt.Sprintf("%08x", len(opHash1)) + hex.EncodeToString(opHash1) +
				fmt.Sprintf("%08x", len(opHash2)) + hex.EncodeToString(opHash2) + "@" +
				hex.EncodeToString(aggregatedSig1))
		execOp1 := []byte(executeBridgeOpsPrefix + "@" +
			hex.EncodeToString(bridgeDataHash1) + "@" +
			hex.EncodeToString([]byte("bridgeOp1")) +
			hex.EncodeToString([]byte("bridgeOp2")))

		registerOp2 := []byte(
			registerBridgeOpsPrefix + "@" +
				hex.EncodeToString(bridgeDataHash2) + "@" +
				fmt.Sprintf("%08x", len(opHash2)) + hex.EncodeToString(opHash3) + "@" +
				hex.EncodeToString(aggregatedSig2))
		execOp2 := []byte(executeBridgeOpsPrefix + "@" +
			hex.EncodeToString(bridgeDataHash2) + "@" +
			hex.EncodeToString([]byte("bridgeOp3")))

		expectedTxsData := [][]byte{
			registerOp1,
			execOp1,
			registerOp2,
			execOp2,
		}

		txsData := df.CreateTxsData(bridgeOps)
		require.Equal(t, expectedTxsData, txsData)
	})
}

package server

import (
	"context"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/data/sovereign/dto"
	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-sovereign-bridge-go/testscommon"
)

func createTxSenders() map[dto.ChainID]TxSender {
	return map[dto.ChainID]TxSender{
		dto.MVX: &testscommon.TxSenderMock{},
		dto.ETH: &testscommon.TxSenderMock{},
	}
}

func TestNewSovereignBridgeTxServer(t *testing.T) {
	t.Parallel()

	t.Run("nil tx sender", func(t *testing.T) {
		bridgeServer, err := NewSovereignBridgeTxServer(nil)
		require.Equal(t, errNilTxSender, err)
		require.Nil(t, bridgeServer)
	})
	t.Run("should work", func(t *testing.T) {
		bridgeServer, err := NewSovereignBridgeTxServer(createTxSenders())
		require.Nil(t, err)
		require.False(t, bridgeServer.IsInterfaceNil())
	})
}

func TestServer_Send(t *testing.T) {
	t.Parallel()

	expectedTxHashes := []string{"txHash1", "txHash2", "txHash3"}
	expectedBridgeOps := &sovereign.BridgeOperations{
		Data: []*sovereign.BridgeOutGoingData{
			{
				ChainID: int32(dto.MVX),
				Hash:    []byte("hash1"),
			},
			{
				ChainID: int32(dto.MVX),
				Hash:    []byte("hash2"),
			},
			{
				ChainID: int32(dto.ETH),
				Hash:    []byte("hash3"),
			},
		},
	}

	txSenderMVX := &testscommon.TxSenderMock{
		SendTxsCalled: func(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
			require.Equal(t, &sovereign.BridgeOperations{
				Data: []*sovereign.BridgeOutGoingData{
					{
						ChainID: int32(dto.MVX),
						Hash:    []byte("hash1"),
					},
					{
						ChainID: int32(dto.MVX),
						Hash:    []byte("hash2"),
					},
				},
			}, data)
			return expectedTxHashes[:2], nil
		},
	}
	txSenderETH := &testscommon.TxSenderMock{
		SendTxsCalled: func(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
			require.Equal(t, &sovereign.BridgeOperations{
				Data: []*sovereign.BridgeOutGoingData{
					{
						ChainID: int32(dto.ETH),
						Hash:    []byte("hash3"),
					},
				},
			}, data)
			return expectedTxHashes[2:], nil
		},
	}

	txsSenders := map[dto.ChainID]TxSender{
		dto.MVX: txSenderMVX,
		dto.ETH: txSenderETH,
	}

	bridgeServer, _ := NewSovereignBridgeTxServer(txsSenders)
	res, err := bridgeServer.Send(context.Background(), expectedBridgeOps)
	require.Nil(t, err)
	require.Equal(t, &sovereign.BridgeOperationsResponse{
		TxHashes: expectedTxHashes,
	}, res)
}

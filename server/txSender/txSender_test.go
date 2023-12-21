package txSender

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/testscommon"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/stretchr/testify/require"
)

func createArgs() TxSenderArgs {
	return TxSenderArgs{
		Wallet:          &testscommon.CryptoComponentsHolderMock{},
		Proxy:           &testscommon.ProxyMock{},
		TxInteractor:    &testscommon.TxInteractorMock{},
		DataFormatter:   &testscommon.DataFormatterMock{},
		SCBridgeAddress: "erd1qqq",
	}
}

func TestNewTxSender(t *testing.T) {
	t.Parallel()

	t.Run("nil wallet", func(t *testing.T) {
		args := createArgs()
		args.Wallet = nil

		ts, err := NewTxSender(args)
		require.Nil(t, ts)
		require.Equal(t, errNilWallet, err)
	})
	t.Run("nil proxy", func(t *testing.T) {
		args := createArgs()
		args.Proxy = nil

		ts, err := NewTxSender(args)
		require.Nil(t, ts)
		require.Equal(t, errNilProxy, err)
	})
	t.Run("nil tx interactor", func(t *testing.T) {
		args := createArgs()
		args.TxInteractor = nil

		ts, err := NewTxSender(args)
		require.Nil(t, ts)
		require.Equal(t, errNilTxInteractor, err)
	})
	t.Run("nil data formatter", func(t *testing.T) {
		args := createArgs()
		args.DataFormatter = nil

		ts, err := NewTxSender(args)
		require.Nil(t, ts)
		require.Equal(t, errNilDataFormatter, err)
	})
	t.Run("should work", func(t *testing.T) {
		args := createArgs()

		ts, err := NewTxSender(args)
		require.Nil(t, err)
		require.False(t, ts.IsInterfaceNil())
	})
}

func TestTxSender_SendTxs(t *testing.T) {
	t.Parallel()

	expectedCtx := context.Background()
	expectedNonce := 0
	expectedTxHashes := []string{"txHash1", "txHash2", "txHash3"}
	expectedTxsData := [][]byte{[]byte("txData1"), []byte("txData2"), []byte("txData3")}
	expectedSigs := []string{"sig1", "sig2", "sig3"}
	expectedBridgeData := &sovereign.BridgeOperations{
		Data: []*sovereign.BridgeOutGoingData{
			{
				Hash: []byte("bridgeDataHash1"),
			},
		},
	}
	expectedNetworkConfig := &data.NetworkConfig{
		MinGasPrice:           5000,
		ChainID:               "1",
		MinTransactionVersion: 2,
	}

	args := createArgs()
	args.Proxy = &testscommon.ProxyMock{
		GetNetworkConfigCalled: func(ctx context.Context) (*data.NetworkConfig, error) {
			require.Equal(t, expectedCtx, ctx)
			return expectedNetworkConfig, nil
		},
	}
	args.DataFormatter = &testscommon.DataFormatterMock{
		CreateTxsDataCalled: func(data *sovereign.BridgeOperations) [][]byte {
			require.Equal(t, expectedBridgeData, data)
			return expectedTxsData
		},
	}
	args.TxInteractor = &testscommon.TxInteractorMock{
		ApplySignatureCalled: func(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
			tx.Signature = expectedSigs[expectedNonce]
			return nil
		},
		AddTransactionCalled: func(tx *transaction.FrontendTransaction) {
			require.Equal(t, &transaction.FrontendTransaction{
				Nonce:     uint64(expectedNonce),
				Value:     "1",
				Receiver:  args.SCBridgeAddress,
				Sender:    args.Wallet.GetBech32(),
				GasPrice:  expectedNetworkConfig.MinGasPrice,
				GasLimit:  50_000_000,
				Data:      expectedTxsData[expectedNonce],
				Signature: expectedSigs[expectedNonce],
				ChainID:   expectedNetworkConfig.ChainID,
				Version:   expectedNetworkConfig.MinTransactionVersion,
			}, tx)

			expectedNonce++
		},
		SendTransactionsAsBunchCalled: func(ctx context.Context, bunchSize int) ([]string, error) {
			require.Equal(t, expectedCtx, ctx)
			require.Equal(t, len(expectedTxHashes), bunchSize)

			return expectedTxHashes, nil
		},
	}

	ts, _ := NewTxSender(args)
	txHashes, err := ts.SendTxs(expectedCtx, expectedBridgeData)
	require.Nil(t, err)
	require.Equal(t, expectedTxHashes, txHashes)
	require.Equal(t, 3, expectedNonce)
}

func TestTxSender_SendTxsConcurrently(t *testing.T) {
	t.Parallel()

	args := createArgs()

	expectedTxHashes := []string{"hash"}
	numTxsToSend := 1000
	numSentTxs := 0
	wg := sync.WaitGroup{}
	wg.Add(numTxsToSend)

	args.Proxy = &testscommon.ProxyMock{
		GetAccountCalled: func(ctx context.Context, address core.AddressHandler) (*data.Account, error) {
			return &data.Account{
				Nonce: uint64(numSentTxs),
			}, nil
		},
	}
	args.DataFormatter = &testscommon.DataFormatterMock{
		CreateTxsDataCalled: func(data *sovereign.BridgeOperations) [][]byte {
			return [][]byte{[]byte("txData")}
		},
	}
	args.TxInteractor = &testscommon.TxInteractorMock{
		SendTransactionsAsBunchCalled: func(ctx context.Context, bunchSize int) ([]string, error) {
			defer func() {
				numSentTxs++
				wg.Done()
			}()

			require.Equal(t, 1, bunchSize)
			return []string{"hash"}, nil
		},
	}

	ts, _ := NewTxSender(args)

	for i := 0; i < numTxsToSend; i++ {
		go func(idx int) {
			txHashes, err := ts.SendTxs(context.Background(), &sovereign.BridgeOperations{
				Data: []*sovereign.BridgeOutGoingData{
					{
						Hash: []byte(fmt.Sprintf("hash%d", idx)),
					},
				},
			})
			require.Nil(t, err)
			require.Equal(t, expectedTxHashes, txHashes)
		}(i)
	}

	wg.Wait()
	require.Equal(t, numTxsToSend, numSentTxs)
}

package txSender

import (
	"testing"

	"github.com/multiversx/mx-chain-sovereign-bridge-go/testscommon"
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

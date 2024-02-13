package deploy

import (
	"testing"

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
		require.Empty(t, df.CreateTxsData(nil))
	})

	t.Run("should work", func(t *testing.T) {
		expectedTxsData := [][]byte{
			{0x01, 0x02},
			VmTypeWasmVm,
			CodeMetadata,
		}

		txsData := df.CreateTxsData([]byte{0x01, 0x02})
		require.Equal(t, expectedTxsData, txsData)
	})
}

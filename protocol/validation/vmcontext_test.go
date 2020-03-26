package validation

import (
	"fmt"
	"testing"

	"github.com/bytom/vapor/errors"
	"github.com/bytom/vapor/protocol/bc"
	"github.com/bytom/vapor/protocol/bc/types"
	"github.com/bytom/vapor/protocol/vm"
)

func TestCheckOutput(t *testing.T) {
	tx := types.NewTx(types.TxData{
		Inputs: []*types.TxInput{
			types.NewSpendInput(nil, bc.Hash{}, bc.NewAssetID([32]byte{1}), 5, 1, []byte("spendprog")),
		},
		Outputs: []*types.TxOutput{
			types.NewIntraChainOutput(bc.NewAssetID([32]byte{3}), 8, []byte("wrongprog")),
			types.NewIntraChainOutput(bc.NewAssetID([32]byte{3}), 8, []byte("controlprog")),
			types.NewIntraChainOutput(bc.NewAssetID([32]byte{2}), 8, []byte("controlprog")),
			types.NewIntraChainOutput(bc.NewAssetID([32]byte{2}), 7, []byte("controlprog")),
			types.NewIntraChainOutput(bc.NewAssetID([32]byte{2}), 7, []byte("controlprog")),
		},
	})

	txCtx := &entryContext{
		entry:   tx.Tx.Entries[tx.Tx.InputIDs[0]],
		entries: tx.Tx.Entries,
	}

	cases := []struct {
		// args to CheckOutput
		index     uint64
		amount    uint64
		assetID   []byte
		vmVersion uint64
		code      []byte

		wantErr error
		wantOk  bool
	}{
		{
			index:     4,
			amount:    7,
			assetID:   append([]byte{2}, make([]byte, 31)...),
			vmVersion: 1,
			code:      []byte("controlprog"),
			wantOk:    true,
		},
		{
			index:     3,
			amount:    7,
			assetID:   append([]byte{2}, make([]byte, 31)...),
			vmVersion: 1,
			code:      []byte("controlprog"),
			wantOk:    true,
		},
		{
			index:     0,
			amount:    1,
			assetID:   append([]byte{9}, make([]byte, 31)...),
			vmVersion: 1,
			code:      []byte("missingprog"),
			wantOk:    false,
		},
		{
			index:     5,
			amount:    7,
			assetID:   append([]byte{2}, make([]byte, 31)...),
			vmVersion: 1,
			code:      []byte("controlprog"),
			wantErr:   vm.ErrBadValue,
		},
	}

	for i, test := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			gotOk, err := txCtx.checkOutput(test.index, test.amount, test.assetID, test.vmVersion, test.code, false)
			if g := errors.Root(err); g != test.wantErr {
				t.Errorf("checkOutput(%v, %v, %x, %v, %x) err = %v, want %v", test.index, test.amount, test.assetID, test.vmVersion, test.code, g, test.wantErr)
				return
			}
			if gotOk != test.wantOk {
				t.Errorf("checkOutput(%v, %v, %x, %v, %x) ok = %t, want %v", test.index, test.amount, test.assetID, test.vmVersion, test.code, gotOk, test.wantOk)
			}

		})
	}
}

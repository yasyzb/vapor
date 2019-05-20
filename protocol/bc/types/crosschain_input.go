package types

import (
	"github.com/vapor/protocol/bc"
)

// CrossChainInput satisfies the TypedInput interface and represents a cross-chain transaction.
type CrossChainInput struct {
	AssetDefinition       []byte
	SpendCommitmentSuffix []byte   // The unconsumed suffix of the spend commitment
	Arguments             [][]byte // Witness
	SpendCommitment
}

// NewCrossChainInput create a new CrossChainInput struct.
// The source is created/issued by trusted federation and hence there is no need
// to refer to it.
func NewCrossChainInput(arguments [][]byte, sourceID bc.Hash, assetID bc.AssetID, amount, sourcePos uint64, controlProgram, assetDefinition []byte) *TxInput {
	sc := SpendCommitment{
		AssetAmount: bc.AssetAmount{
			AssetId: &assetID,
			Amount:  amount,
		},
		SourceID:       sourceID,
		SourcePosition: sourcePos,
		VMVersion:      1,
		ControlProgram: controlProgram,
	}
	return &TxInput{
		AssetVersion: 1,
		TypedInput: &CrossChainInput{
			AssetDefinition: assetDefinition,
			SpendCommitment: sc,
			Arguments:       arguments,
		},
	}
}

// InputType is the interface function for return the input type.
func (si *CrossChainInput) InputType() uint8 { return CrossChainInputType }
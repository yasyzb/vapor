package consensus

import (
	"github.com/vapor/chain"
	"github.com/vapor/protocol/bc"
	"github.com/vapor/protocol/bc/types"
)

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	Author(header *types.BlockHeader, c chain.Chain) (string, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(c chain.Chain, header *types.BlockHeader, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(c chain.Chain, headers []*types.BlockHeader, seals []bool) (chan<- struct{}, <-chan error)

	// VerifySeal checks whether the crypto seal on a header is valid according to
	// the consensus rules of the given engine.
	VerifySeal(c chain.Chain, header *types.BlockHeader) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(c chain.Chain, header *types.BlockHeader) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// and assembles the final block.
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(c chain.Chain, header *types.BlockHeader, txs []*bc.Tx) error

	// Seal generates a new block for the given input block with the local miner's
	// seal place on top.
	Seal(c chain.Chain, block *types.Block) (*types.Block, error)
}

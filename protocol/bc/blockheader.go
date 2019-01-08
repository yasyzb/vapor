package bc

import "io"

// BlockHeader contains the header information for a blockchain
// block. It satisfies the Entry interface.

func (BlockHeader) typ() string { return "blockheader" }
func (bh *BlockHeader) writeForHash(w io.Writer) {
	mustWriteForHash(w, bh.Version)
	mustWriteForHash(w, bh.Height)
	mustWriteForHash(w, bh.PreviousBlockId)
	mustWriteForHash(w, bh.Timestamp)
	//mustWriteForHash(w, bh.Coinbase)
	mustWriteForHash(w, bh.TransactionsRoot)
	mustWriteForHash(w, bh.TransactionStatusHash)
	mustWriteForHash(w, bh.Proof.Sign)
	mustWriteForHash(w, bh.Proof.ControlProgram)
	mustWriteForHash(w, bh.Extra)
}

// NewBlockHeader creates a new BlockHeader and populates
// its body.
func NewBlockHeader(version, height uint64, previousBlockID *Hash, timestamp uint64, transactionsRoot, transactionStatusHash *Hash, proof *Proof, extra []byte, coinbase []byte) *BlockHeader {
	return &BlockHeader{
		Version:               version,
		Height:                height,
		PreviousBlockId:       previousBlockID,
		Timestamp:             timestamp,
		TransactionsRoot:      transactionsRoot,
		TransactionStatusHash: transactionStatusHash,
		TransactionStatus:     nil,
		Proof:                 proof,
		Extra:                 extra,
		Coinbase:              coinbase,
	}
}

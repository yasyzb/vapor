package api

import (
	set "gopkg.in/fatih/set.v0"

	"github.com/vapor/blockchain/query"
	"github.com/vapor/crypto/ed25519/chainkd"
	chainjson "github.com/vapor/encoding/json"
	"github.com/vapor/protocol/bc"
	"github.com/vapor/protocol/bc/types"
	"github.com/vapor/protocol/state"
)

// return best block hash
func (a *API) getBestBlockHash() Response {
	blockHash := map[string]string{"block_hash": a.chain.BestBlockHash().String()}
	return NewSuccessResponse(blockHash)
}

// return current block count
func (a *API) getBlockCount() Response {
	blockHeight := map[string]uint64{"block_count": a.chain.BestBlockHeight()}
	return NewSuccessResponse(blockHeight)
}

// BlockTx is the tx struct for getBlock func
type BlockTx struct {
	ID         bc.Hash                  `json:"id"`
	Version    uint64                   `json:"version"`
	Size       uint64                   `json:"size"`
	TimeRange  uint64                   `json:"time_range"`
	Inputs     []*query.AnnotatedInput  `json:"inputs"`
	Outputs    []*query.AnnotatedOutput `json:"outputs"`
	StatusFail bool                     `json:"status_fail"`
	MuxID      bc.Hash                  `json:"mux_id"`
}

// BlockReq is used to handle getBlock req
type BlockReq struct {
	BlockHeight uint64             `json:"block_height"`
	BlockHash   chainjson.HexBytes `json:"block_hash"`
}

// GetBlockResp is the resp for getBlock api
type GetBlockResp struct {
	Hash                   *bc.Hash             `json:"hash"`
	Size                   uint64               `json:"size"`
	Version                uint64               `json:"version"`
	Height                 uint64               `json:"height"`
	PreviousBlockHash      *bc.Hash             `json:"previous_block_hash"`
	Timestamp              uint64               `json:"timestamp"`
	Witness                []chainjson.HexBytes `json:"witness"`
	TransactionsMerkleRoot *bc.Hash             `json:"transaction_merkle_root"`
	TransactionStatusHash  *bc.Hash             `json:"transaction_status_hash"`
	Transactions           []*BlockTx           `json:"transactions"`
}

// return block by hash/height
func (a *API) getBlock(ins BlockReq) Response {
	block, err := a.getBlockHelper(ins)
	if err != nil {
		return NewErrorResponse(err)
	}

	blockHash := block.Hash()
	txStatus, err := a.chain.GetTransactionStatus(&blockHash)
	rawBlock, err := block.MarshalText()
	if err != nil {
		return NewErrorResponse(err)
	}

	witness := make([]chainjson.HexBytes, len(block.Witness))
	for i, w := range block.Witness {
		witness[i] = w
	}

	resp := &GetBlockResp{
		Hash:                   &blockHash,
		Size:                   uint64(len(rawBlock)),
		Version:                block.Version,
		Height:                 block.Height,
		PreviousBlockHash:      &block.PreviousBlockHash,
		Timestamp:              block.Timestamp,
		Witness:                witness,
		TransactionsMerkleRoot: &block.TransactionsMerkleRoot,
		TransactionStatusHash:  &block.TransactionStatusHash,
		Transactions:           []*BlockTx{},
	}

	for i, orig := range block.Transactions {
		tx := &BlockTx{
			ID:        orig.ID,
			Version:   orig.Version,
			Size:      orig.SerializedSize,
			TimeRange: orig.TimeRange,
			Inputs:    []*query.AnnotatedInput{},
			Outputs:   []*query.AnnotatedOutput{},
		}
		tx.StatusFail, err = txStatus.GetStatus(i)
		if err != nil {
			return NewSuccessResponse(resp)
		}

		resOutID := orig.ResultIds[0]
		switch resOut := orig.Entries[*resOutID].(type) {
		case *bc.IntraChainOutput:
			tx.MuxID = *resOut.Source.Ref
		case *bc.CrossChainOutput:
			tx.MuxID = *resOut.Source.Ref
		case *bc.VoteOutput:
			tx.MuxID = *resOut.Source.Ref
		case *bc.Retirement:
			resRetire, _ := orig.Entries[*resOutID].(*bc.Retirement)
			tx.MuxID = *resRetire.Source.Ref
		}

		for i := range orig.Inputs {
			tx.Inputs = append(tx.Inputs, a.wallet.BuildAnnotatedInput(orig, uint32(i)))
		}
		for i := range orig.Outputs {
			tx.Outputs = append(tx.Outputs, a.wallet.BuildAnnotatedOutput(orig, i))
		}
		resp.Transactions = append(resp.Transactions, tx)
	}
	return NewSuccessResponse(resp)
}

// GetRawBlockResp is resp struct for getRawBlock API
type GetRawBlockResp struct {
	RawBlock          *types.Block          `json:"raw_block"`
	TransactionStatus *bc.TransactionStatus `json:"transaction_status"`
}

func (a *API) getRawBlock(ins BlockReq) Response {
	block, err := a.getBlockHelper(ins)
	if err != nil {
		return NewErrorResponse(err)
	}

	blockHash := block.Hash()
	txStatus, err := a.chain.GetTransactionStatus(&blockHash)
	if err != nil {
		return NewErrorResponse(err)
	}

	resp := GetRawBlockResp{
		RawBlock:          block,
		TransactionStatus: txStatus,
	}
	return NewSuccessResponse(resp)
}

// GetBlockHeaderResp is resp struct for getBlockHeader API
type GetBlockHeaderResp struct {
	BlockHeader *types.BlockHeader `json:"block_header"`
	Reward      uint64             `json:"reward"`
}

func (a *API) getBlockHeader(ins BlockReq) Response {
	block, err := a.getBlockHelper(ins)
	if err != nil {
		return NewErrorResponse(err)
	}

	resp := &GetBlockHeaderResp{
		BlockHeader: &block.BlockHeader,
		Reward:      block.Transactions[0].Outputs[0].AssetAmount().Amount,
	}
	return NewSuccessResponse(resp)
}

func (a *API) getBlockHelper(ins BlockReq) (*types.Block, error) {
	if len(ins.BlockHash) == 32 {
		hash := hexBytesToHash(ins.BlockHash)
		return a.chain.GetBlockByHash(&hash)
	} else {
		return a.chain.GetBlockByHeight(ins.BlockHeight)
	}
}

func hexBytesToHash(hexBytes chainjson.HexBytes) bc.Hash {
	b32 := [32]byte{}
	copy(b32[:], hexBytes)
	return bc.NewHash(b32)
}

// MerkleBlockReq is used to handle getTxOutProof req
type MerkleBlockReq struct {
	TxIDs     []chainjson.HexBytes `json:"tx_ids"`
	BlockHash chainjson.HexBytes   `json:"block_hash"`
}

// GetMerkleBlockResp is resp struct for GetTxOutProof API
type GetMerkleBlockResp struct {
	BlockHeader  types.BlockHeader `json:"block_header"`
	TxHashes     []*bc.Hash        `json:"tx_hashes"`
	StatusHashes []*bc.Hash        `json:"status_hashes"`
	Flags        []uint32          `json:"flags"`
	MatchedTxIDs []*bc.Hash        `json:"matched_tx_ids"`
}

func (a *API) getMerkleProof(ins MerkleBlockReq) Response {
	blockReq := BlockReq{BlockHash: ins.BlockHash}
	block, err := a.getBlockHelper(blockReq)
	if err != nil {
		return NewErrorResponse(err)
	}

	matchedTxs := getMatchedTx(block.Transactions, ins.TxIDs)
	var matchedTxIDs []*bc.Hash
	for _, tx := range matchedTxs {
		matchedTxIDs = append(matchedTxIDs, &tx.ID)
	}

	hashes, compactFlags := types.GetTxMerkleTreeProof(block.Transactions, matchedTxs)
	flags := make([]uint32, len(compactFlags))
	for i, flag := range compactFlags {
		flags[i] = uint32(flag)
	}

	blockHash := block.Hash()
	statuses, err := a.chain.GetTransactionStatus(&blockHash)
	if err != nil {
		return NewErrorResponse(err)
	}

	statusHashes := types.GetStatusMerkleTreeProof(statuses.VerifyStatus, compactFlags)

	resp := &GetMerkleBlockResp{
		BlockHeader:  block.BlockHeader,
		TxHashes:     hashes,
		StatusHashes: statusHashes,
		Flags:        flags,
		MatchedTxIDs: matchedTxIDs,
	}
	return NewSuccessResponse(resp)
}

func getMatchedTx(txs []*types.Tx, filterTxIDs []chainjson.HexBytes) []*types.Tx {
	txIDSet := set.New()
	for _, txID := range filterTxIDs {
		hash := hexBytesToHash(txID)
		txIDSet.Add(hash.String())
	}

	var matchedTxs []*types.Tx
	for _, tx := range txs {
		hashStr := tx.ID.String()
		if txIDSet.Has(hashStr) {
			matchedTxs = append(matchedTxs, tx)
		}
	}
	return matchedTxs
}

// GetVoteResultResp is used to handle GetVoteResult req
type GetVoteResultResp struct {
	XPub    chainkd.XPub `json:"xPub"`
	VoteNum uint64       `json:"vote_num"`
	Order   uint64       `json:"order"`
}

// GetVoteResult return VoteResult by hash/height
func (a *API) GetVoteResult(ins BlockReq) Response {
	resp := []GetVoteResultResp{}
	voteResult := make(map[string]*state.ConsensusNode)
	var err error

	if len(ins.BlockHash) == 32 {
		hash := hexBytesToHash(ins.BlockHash)
		voteResult, err = a.chain.GetVoteResultByHash(&hash)
	} else {
		voteResult, err = a.chain.GetVoteResultByHeight(ins.BlockHeight)
	}
	if err != nil {
		return NewErrorResponse(err)
	}
	for _, v := range voteResult {
		resp = append(resp, GetVoteResultResp{
			XPub:    v.XPub,
			VoteNum: v.VoteNum,
			Order:   v.Order,
		})
	}

	return NewSuccessResponse(resp)
}

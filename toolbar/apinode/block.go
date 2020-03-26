package apinode

import (
	"encoding/json"

	"github.com/bytom/vapor/api"
	"github.com/bytom/vapor/errors"
	"github.com/bytom/vapor/protocol/bc/types"
)

func (n *Node) GetBlockByHash(hash string) (*types.Block, error) {
	return n.getRawBlock(&getRawBlockReq{BlockHash: hash})
}

func (n *Node) GetBlockByHeight(height uint64) (*types.Block, error) {
	return n.getRawBlock(&getRawBlockReq{BlockHeight: height})
}

type getRawBlockReq struct {
	BlockHeight uint64 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
}

func (n *Node) getRawBlock(req *getRawBlockReq) (*types.Block, error) {
	url := "/get-raw-block"
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}
	resp := &api.GetRawBlockResp{}
	return resp.RawBlock, n.request(url, payload, resp)
}

func (n *Node) GetVoteByHash(hash string) ([]*api.VoteInfo, error) {
	return n.getVoteResult(&getVoteResultReq{BlockHash: hash})
}

func (n *Node) GetVoteByHeight(height uint64) ([]*api.VoteInfo, error) {
	return n.getVoteResult(&getVoteResultReq{BlockHeight: height})
}

type getVoteResultReq struct {
	BlockHeight uint64 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
}

func (n *Node) getVoteResult(req *getVoteResultReq) ([]*api.VoteInfo, error) {
	url := "/get-vote-result"
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}
	resp := []*api.VoteInfo{}
	return resp, n.request(url, payload, &resp)
}

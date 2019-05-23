package protocol

import (
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/vapor/errors"
	"github.com/vapor/protocol/state"
)

const (
	numOfConsensusNode = 21
	roundVoteBlockNums = 1000

	// product one block per 500 milliseconds
	blockTimeInterval = 500
	blockNumEachNode  = 3
)

var (
	errHasNoChanceProductBlock  = errors.New("the node has no chance to product a block in this round of voting")
	errNotFoundConsensusNode    = errors.New("can not found consensus node")
	errVoteResultIsNotfinalized = errors.New("vote result is not finalized")
)

type consensusNode struct {
	pubkey  string
	voteNum uint64
	order   uint64
}

type consensusNodeSlice []*consensusNode

func (c consensusNodeSlice) Len() int           { return len(c) }
func (c consensusNodeSlice) Less(i, j int) bool { return c[i].voteNum > c[j].voteNum }
func (c consensusNodeSlice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type consensusNodeManager struct {
	consensusNodeMap     map[string]*consensusNode
	effectiveStartHeight uint64
	store                Store
	blockIndex           *state.BlockIndex
	sync.RWMutex
}

func newConsensusNodeManager(store Store, blockIndex *state.BlockIndex) *consensusNodeManager {
	return &consensusNodeManager{
		consensusNodeMap:     make(map[string]*consensusNode),
		effectiveStartHeight: 1,
		store:                store,
		blockIndex:           blockIndex,
	}
}

func (c *consensusNodeManager) getConsensusNode(height uint64, pubkey string) (*consensusNode, error) {
	defer c.RUnlock()
	c.RLock()
	if height >= c.effectiveStartHeight+roundVoteBlockNums {
		return nil, errors.New("the vote has not been completed for the specified block height ")
	}

	var err error
	consensusNodeMap := c.consensusNodeMap
	// query history vote result
	if height < c.effectiveStartHeight {
		consensusNodeMap, err = c.getConsensusNodesByVoteResult(height)
		if err != nil {
			return nil, err
		}
	}

	node, exist := consensusNodeMap[pubkey]
	if !exist {
		return node, errNotFoundConsensusNode
	}
	return node, nil
}

func (c *consensusNodeManager) isBlocker(height uint64, blockTimestamp uint64, pubkey string) (bool, error) {
	prevVoteRoundLastBlock := c.blockIndex.NodeByHeight(height - 1)
	startTimestamp := prevVoteRoundLastBlock.Timestamp + blockTimeInterval

	consensusNodeMap, err := c.getConsensusNodesByVoteResult(height)
	if err != nil {
		return false, err
	}

	blockerNode, exist := consensusNodeMap[pubkey]
	if !exist {
		return false, nil
	}

	begin := getLastBlockTimeInTimeRange(startTimestamp, blockTimestamp, blockerNode.order)
	end := begin + blockNumEachNode*blockTimeInterval
	return blockTimestamp >= begin && blockTimestamp < end, nil
}

func (c *consensusNodeManager) nextLeaderTime(pubkey []byte, bestBlockTimestamp, bestBlockHeight uint64) (*time.Time, error) {
	defer c.RUnlock()
	c.RLock()

	startHeight := c.effectiveStartHeight
	prevRoundLastBlock := c.blockIndex.NodeByHeight(startHeight - 1)
	startTime := prevRoundLastBlock.Timestamp + blockTimeInterval
	endTime := bestBlockTimestamp + (roundVoteBlockNums-bestBlockHeight%roundVoteBlockNums)*blockTimeInterval

	consensusNode, exist := c.consensusNodeMap[hex.EncodeToString(pubkey)]
	if !exist {
		return nil, fmt.Errorf("pubkey:%s is not consensus node", hex.EncodeToString(pubkey))
	}

	nextLeaderTime, err := nextLeaderTimeHelper(startTime, endTime, uint64(time.Now().UnixNano()/1e6), consensusNode.order)
	if err != nil {
		return nil, err
	}

	return nextLeaderTime, nil
}

func nextLeaderTimeHelper(startTime, endTime, now, nodeOrder uint64) (*time.Time, error) {
	nextLeaderTimestamp := getLastBlockTimeInTimeRange(startTime, now, nodeOrder)
	roundBlockTime := uint64(blockNumEachNode * numOfConsensusNode * blockTimeInterval)

	if int64(now-nextLeaderTimestamp) >= blockNumEachNode*blockTimeInterval {
		nextLeaderTimestamp += roundBlockTime
		if nextLeaderTimestamp >= endTime {
			return nil, errHasNoChanceProductBlock
		}
	}

	nextLeaderTime := time.Unix(int64(nextLeaderTimestamp)/1000, (int64(nextLeaderTimestamp)%1000)*1e6)
	return &nextLeaderTime, nil
}

// updateConsensusNodes used to update consensus node after each round of voting
func (c *consensusNodeManager) updateConsensusNodes(bestBlockHeight uint64) error {
	defer c.Unlock()
	c.Lock()

	consensusNodeMap, err := c.getConsensusNodesByVoteResult(bestBlockHeight)
	if err != nil && err != errVoteResultIsNotfinalized {
		return err
	}

	if err == errVoteResultIsNotfinalized {
		return nil
	}

	c.consensusNodeMap = consensusNodeMap
	c.effectiveStartHeight = bestBlockHeight / roundVoteBlockNums * roundVoteBlockNums
	return nil
}

func getLastBlockTimeInTimeRange(startTimestamp, endTimestamp, order uint64) uint64 {
	// One round of product block time for all consensus nodes
	roundBlockTime := uint64(blockNumEachNode * numOfConsensusNode * blockTimeInterval)
	// The start time of the last round of product block
	lastRoundStartTime := startTimestamp + (endTimestamp-startTimestamp)/roundBlockTime*roundBlockTime
	// The time of product block of the consensus in last round
	return lastRoundStartTime + order*(blockNumEachNode*blockTimeInterval)
}

func (c *consensusNodeManager) getConsensusNodesByVoteResult(blockHeight uint64) (map[string]*consensusNode, error) {
	defer c.RUnlock()
	c.RLock()
	if blockHeight >= c.effectiveStartHeight+roundVoteBlockNums {
		return nil, errors.New("the given block height is greater than current vote start height")
	}

	if blockHeight >= c.effectiveStartHeight {
		return c.consensusNodeMap, nil
	}

	voteResult, err := c.store.GetVoteResult(blockHeight / roundVoteBlockNums)
	if err != nil {
		return nil, err
	}

	if !voteResult.Finalized {
		return nil, errVoteResultIsNotfinalized
	}

	var nodes []*consensusNode
	for pubkey, voteNum := range voteResult.NumOfVote {
		nodes = append(nodes, &consensusNode{
			pubkey:  pubkey,
			voteNum: voteNum,
		})
	}
	// In principle, there is no need to sort all voting nodes.
	// if there is a performance problem, consider the optimization later.
	// TODO not consider the same number of votes
	sort.Sort(consensusNodeSlice(nodes))

	result := make(map[string]*consensusNode)
	for i := 0; i < numOfConsensusNode; i++ {
		node := nodes[i]
		node.order = uint64(i)
		result[node.pubkey] = node
	}
	return result, nil
}
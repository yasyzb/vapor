package chainmgr

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/vapor/consensus"
	"github.com/vapor/errors"
	"github.com/vapor/netsync/peers"
	"github.com/vapor/p2p/security"
	"github.com/vapor/protocol/bc"
	"github.com/vapor/protocol/bc/types"
)

const (
	syncCycle            = 5 * time.Second
	blockProcessChSize   = 1024
	blocksProcessChSize  = 128
	headersProcessChSize = 1024
)

var (
	syncTimeout = 30 * time.Second

	errRequestTimeout = errors.New("request timeout")
	errPeerDropped    = errors.New("Peer dropped")
)

type FastSync interface {
	locateBlocks(locator []*bc.Hash, stopHash *bc.Hash) ([]*types.Block, error)
	locateHeaders(locator []*bc.Hash, stopHash *bc.Hash, amount uint64, skip uint64, maxNum uint64) ([]*types.BlockHeader, error)
	process() error
	processBlocks(peerID string, blocks []*types.Block)
	processHeaders(peerID string, headers []*types.BlockHeader)
	//setSyncPeer(peer *peers.Peer)
	stop()
}

type blockMsg struct {
	block  *types.Block
	peerID string
}

type blocksMsg struct {
	blocks []*types.Block
	peerID string
}

type headersMsg struct {
	headers []*types.BlockHeader
	peerID  string
}

type blockKeeper struct {
	chain    Chain
	fastSync FastSync
	peers    *peers.PeerSet

	syncPeer       *peers.Peer
	blockProcessCh chan *blockMsg
}

func newBlockKeeper(chain Chain, peers *peers.PeerSet) *blockKeeper {
	return &blockKeeper{
		chain:          chain,
		fastSync:       newFastSync(chain, peers),
		peers:          peers,
		blockProcessCh: make(chan *blockMsg, blockProcessChSize),
	}
}

func (bk *blockKeeper) locateBlocks(locator []*bc.Hash, stopHash *bc.Hash) ([]*types.Block, error) {
	return bk.fastSync.locateBlocks(locator, stopHash)
}

func (bk *blockKeeper) locateHeaders(locator []*bc.Hash, stopHash *bc.Hash, amount uint64, skip uint64, maxNum uint64) ([]*types.BlockHeader, error) {
	return bk.fastSync.locateHeaders(locator, stopHash, amount, skip, maxNum)
}

func (bk *blockKeeper) processBlock(peerID string, block *types.Block) {
	bk.blockProcessCh <- &blockMsg{block: block, peerID: peerID}
}

func (bk *blockKeeper) processBlocks(peerID string, blocks []*types.Block) {
	bk.fastSync.processBlocks(peerID, blocks)
}

func (bk *blockKeeper) processHeaders(peerID string, headers []*types.BlockHeader) {
	bk.fastSync.processHeaders(peerID, headers)
}

func (bk *blockKeeper) regularBlockSync() error {
	peer := bk.peers.BestPeer(consensus.SFFastSync | consensus.SFFullNode)
	if peer == nil {
		log.WithFields(log.Fields{"module": logModule}).Debug("can't find sync peer")
		return nil
	}
	peerHeight := peer.Height()
	if peerHeight >= bk.chain.BestBlockHeight()+minGapStartFastSync {
		log.WithFields(log.Fields{"module": logModule}).Debug("Height gap meet fast synchronization condition")
		return nil
	}

	i := bk.chain.BestBlockHeight() + 1
	for i <= peerHeight {
		block, err := bk.requireBlock(i)
		if err != nil {
			bk.peers.ErrorHandler(peer.ID(), security.LevelConnException, err)
			return err
		}

		isOrphan, err := bk.chain.ProcessBlock(block)
		if err != nil {
			bk.peers.ErrorHandler(peer.ID(), security.LevelMsgIllegal, err)
			return err
		}

		if isOrphan {
			i--
			continue
		}
		i = bk.chain.BestBlockHeight() + 1
	}
	return nil
}

func (bk *blockKeeper) requireBlock(height uint64) (*types.Block, error) {
	if ok := bk.syncPeer.GetBlockByHeight(height); !ok {
		return nil, errPeerDropped
	}

	timeout := time.NewTimer(syncTimeout)
	defer timeout.Stop()

	for {
		select {
		case msg := <-bk.blockProcessCh:
			if msg.peerID != bk.syncPeer.ID() {
				continue
			}
			if msg.block.Height != height {
				continue
			}
			return msg.block, nil
		case <-timeout.C:
			return nil, errors.Wrap(errRequestTimeout, "requireBlock")
		}
	}
}

func (bk *blockKeeper) start() {
	go bk.syncWorker()
}

func (bk *blockKeeper) startSync() bool {
	if err := bk.fastSync.process(); err != nil {
		log.WithFields(log.Fields{"module": logModule, "err": err}).Warning("failed on fast sync")
		return false
	}

	if err := bk.regularBlockSync(); err != nil {
		log.WithFields(log.Fields{"module": logModule, "err": err}).Warning("fail on regularBlockSync")
		return false
	}
	return true
}

func (bk *blockKeeper) stop() {
	bk.fastSync.stop()
}

func (bk *blockKeeper) syncWorker() {
	syncTicker := time.NewTicker(syncCycle)
	defer syncTicker.Stop()

	for {
		select {
		case <-syncTicker.C:
			if update := bk.startSync(); !update {
				continue
			}

			block, err := bk.chain.GetBlockByHeight(bk.chain.BestBlockHeight())
			if err != nil {
				log.WithFields(log.Fields{"module": logModule, "err": err}).Error("fail on syncWorker get best block")
			}

			if err = bk.peers.BroadcastNewStatus(block); err != nil {
				log.WithFields(log.Fields{"module": logModule, "err": err}).Error("fail on syncWorker broadcast new status")
			}
		}
	}
}

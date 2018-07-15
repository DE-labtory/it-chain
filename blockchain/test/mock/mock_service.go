package mock

import "github.com/it-chain/it-chain-Engine/blockchain"

type SyncService struct {
	SyncWithPeerFunc func(peer blockchain.Peer) error
	SyncedCheckFunc  func(peer blockchain.Peer) (blockchain.IsSynced, error)
	ConstructFunc    func(peer blockchain.Peer) error
}

func (ss SyncService) SyncWithPeer(peer blockchain.Peer) error {
	return ss.SyncWithPeerFunc(peer)
}

func (ss SyncService) SyncedCheck(peer blockchain.Peer) (blockchain.IsSynced, error) {
	return ss.SyncedCheckFunc(peer)
}

func (ss SyncService) Construct(peer blockchain.Peer) error {
	return ss.ConstructFunc(peer)
}

type BlockQueryService struct {
	GetLastBlockFunc             func() (blockchain.Block, error)
	GetBlockByHeightFunc         func(blockHeight uint64) (blockchain.Block, error)
	GetBlockBySealFunc           func(seal []byte) (blockchain.Block, error)
	GetStagedBlockByHeightFunc   func(blockHeight uint64) (blockchain.Block, error)
	GetStagedBlockByIdFunc       func(blockId string) (blockchain.Block, error)
	GetLastCommitedBlockFunc     func() (blockchain.Block, error)
	GetCommitedBlockByHeightFunc func(blockHeight uint64) (blockchain.Block, error)
	GetLastBlockFromPeerFunc     func(peer blockchain.Peer) (blockchain.Block, error)
	GetBlockByHeightFromPeerFunc func(peer blockchain.Peer, height blockchain.BlockHeight) (blockchain.Block, error)
}

func (bqs BlockQueryService) GetLastBlock() (blockchain.Block, error) {
	return bqs.GetLastBlockFunc()
}
func (bqs BlockQueryService) GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.Block, error) {
	return bqs.GetLastBlockFromPeerFunc(peer)
}

func (bqs BlockQueryService) GetBlockBySeal(seal []byte) (blockchain.Block, error) {
	return bqs.GetBlockBySealFunc(seal)
}

func (bqs BlockQueryService) GetBlockByHeightFromPeer(peer blockchain.Peer, height blockchain.BlockHeight) (blockchain.Block, error) {
	return bqs.GetBlockByHeightFromPeerFunc(peer, height)
}
func (bqs BlockQueryService) GetBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return bqs.GetBlockByHeightFunc(blockHeight)
}
func (bqs BlockQueryService) GetStagedBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return bqs.GetStagedBlockByHeightFunc(blockHeight)
}
func (bqs BlockQueryService) GetStagedBlockById(blockId string) (blockchain.Block, error) {
	return bqs.GetStagedBlockByIdFunc(blockId)
}
func (bqs BlockQueryService) GetLastCommitedBlock() (blockchain.Block, error) {
	return bqs.GetLastCommitedBlockFunc()
}
func (bqs BlockQueryService) GetCommitedBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return bqs.GetCommitedBlockByHeightFunc(blockHeight)
}

type PeerService struct {
	GetRandomPeerFunc func() (blockchain.Peer, error)
}

func (ps PeerService) GetRandomPeer() (blockchain.Peer, error) {
	return ps.GetRandomPeerFunc()
}

type BlockService struct {
	ExecuteBlockFunc func(block blockchain.Block) error
}

func (b BlockService) ExecuteBlock(block blockchain.Block) error {
	return b.ExecuteBlockFunc(block)
}

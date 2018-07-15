package mock

import "github.com/it-chain/it-chain-Engine/blockchain"

type SyncService struct {
	BlockQueryService BlockQueryService
	SyncWithPeerFunc  func(peer blockchain.Peer) error
	SyncedCheckFunc   func(peer blockchain.Peer) (blockchain.IsSynced, error)
	ConstructFunc     func(peer blockchain.Peer) error
}

func (ss SyncService) SyncWithPeer(peer blockchain.Peer) error {
	return ss.SyncWithPeerFunc(peer)
}

func (ss SyncService) syncedCheck(peer blockchain.Peer) (blockchain.IsSynced, error) {
	return ss.SyncedCheckFunc(peer)
}

func (ss SyncService) construct(peer blockchain.Peer) error {
	return ss.ConstructFunc(peer)
}

type BlockQueryService struct {
	GetLastBlockFunc             func() (blockchain.Block, error)
	GetLastBlockFromPeerFunc     func(peer blockchain.Peer) (blockchain.Block, error)
	GetBlockByHeightFromPeerFunc func(peer blockchain.Peer, height blockchain.BlockHeight) (blockchain.Block, error)
}

func (bqs BlockQueryService) GetLastBlock() (blockchain.Block, error) {
	return bqs.GetLastBlockFunc()
}

func (bqs BlockQueryService) GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.Block, error) {
	return bqs.GetLastBlockFromPeerFunc(peer)
}

func (bqs BlockQueryService) GetBlockByHeightFromPeer(peer blockchain.Peer, height blockchain.BlockHeight) (blockchain.Block, error) {
	return bqs.GetBlockByHeightFromPeerFunc(peer, height)
}

type PeerService struct {
	GetRandomPeerFunc func() (blockchain.Peer, error)
}

func (ps PeerService) GetRandomPeer() (blockchain.Peer, error) {
	return ps.GetRandomPeerFunc()
}

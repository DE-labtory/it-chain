package adapter

import "github.com/it-chain/engine/blockchain"

type QuerySerivce struct {
}

func NewQueryService() QuerySerivce {
	return QuerySerivce{}
}

//ToDo: Implementation
func (s QuerySerivce) GetRandomPeer() (blockchain.Peer, error) {
	return blockchain.Peer{}, nil
}

//ToDo: Implementation
func (s QuerySerivce) GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.DefaultBlock, error) {
	return blockchain.DefaultBlock{}, nil
}

//ToDo: Implementation
func (s QuerySerivce) GetBlockByHeightFromPeer(peer blockchain.Peer, height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return blockchain.DefaultBlock{}, nil
}

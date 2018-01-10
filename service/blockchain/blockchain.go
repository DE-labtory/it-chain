package blockchain

import (
	"sync"
)

const (
	defaultChannelName = "0"
	defaultPeerId      = "0"
)

type ChainHeader struct {
	chainHeight int    //height of chain
	channelName string //channel name
	peerID      string //owner peer id of chain
}

type BlockChain struct {
	sync.RWMutex              //lock
	header       *ChainHeader //chain meta information
	blocks       []*Block     //list of bloc
}

func CreateNewBlockChain(channelID string, peerId string) *BlockChain {

	var header = ChainHeader{
		chainHeight: 0,
		channelName: channelID,
		peerID:      peerId,
	}

	return &BlockChain{header: &header, blocks: make([]*Block, 0)}
}
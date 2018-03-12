package domain

import (
	"sync"
)

const (
	defaultChannelName = "0"
	defaultPeerId      = "0"
)

type ChainHeader struct {
	ChainHeight int    //height of chain
	ChannelName string //channel name
	PeerID      string //owner peer id of chain
}

type BlockChain struct {
	sync.RWMutex              //lock
	Header       *ChainHeader //chain meta information
	Blocks       []*Block     //list of bloc
}

func CreateNewBlockChain(channelID string, peerId string) *BlockChain {

	var header = ChainHeader{
		ChainHeight: 0,
		ChannelName: channelID,
		PeerID:      peerId,
	}

	return &BlockChain{Header: &header, Blocks: make([]*Block, 0)}
}
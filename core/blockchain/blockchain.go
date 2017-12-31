package blockchain

import (
	"sync"
	"time"
)

type BlockStatus int

const (
	unconfirmed BlockStatus = 0 + iota //unconfirmed block
	confirmed
)

const(
	defaultChannelName = "0"
	defaultPeerId = "0"
)


type ChainHeader struct{
	chainHeight int 				//height of chain
	channelName string 				//channel name
	peerID string 					//owner peer id of chain
}

type Block struct{
	version string 					//version of block
	previousBlockHash string 		//hash of previous block
	merkleTreeRootHash string
	merkleTree []*Transaction
	timeStamp time.Time
	blockHeight int
	blockStatus BlockStatus
	createdPeerID string
	signature []byte				//
}

type BlockChain struct{
	Mux sync.RWMutex				//lock
	Header *ChainHeader				//chain meta information
	Blocks []*Block					//list of block
}

func CreateNewBlockChain(channelID string,peerId string) *BlockChain{

	var header = ChainHeader{
		chainHeight: 0,
		channelName: channelID,
		peerID: peerId,
	}

	return &BlockChain{Header:&header,Blocks:make([]*Block,0)}
}
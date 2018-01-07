package blockchain

import (
	"sync"
	"time"
)

type BlockStatus int

const (
	blockUnconfirmed BlockStatus = 0 + iota //unconfirmed block
	blockConfirmed
)

const(
	defaultChannelName = "0"
	defaultPeerId = "0"
)

type ChainHeader struct {
	chainHeight int    //height of chain
	channelName string //channel name
	peerID      string //owner peer id of chain
}

type Block struct {
	Version            string //version of block
	PreviousBlockHash  string //hash of previous block
	MerkleTreeRootHash string
	MerkleTree         []*Transaction
	TimeStamp          time.Time
	BlockHeight        int
	BlockStatus        BlockStatus
	CreatedPeerID      string
	Signature          []byte
}

type BlockChain struct {
	sync.RWMutex //lock
	Header *ChainHeader //chain meta information
	Blocks []*Block     //list of bloc
}

func CreateNewBlockChain(channelID string,peerId string) *BlockChain{

	var header = ChainHeader{
		chainHeight: 0,
		channelName: channelID,
		peerID: peerId,
	}

	return &BlockChain{Header:&header,Blocks:make([]*Block,0)}
}
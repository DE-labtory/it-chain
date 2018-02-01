package blockchain

import (
	"sync"
	//"errors"
	"it-chain/db/blockchaindb/blockchainleveldb"
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
	sync.RWMutex           //lock
	Header    *ChainHeader //chain meta information
	Blocks    []*Block     //list of bloc
	LastBlock *Block       //last block address
	BlockNum  int          //number of block
	DB        blockchainleveldb.BlockchainLevelDB
}

func CreateNewBlockChain(channelID string, peerId string) *BlockChain {
	var header = ChainHeader{
		ChainHeight: 0,
		ChannelName: channelID,
		PeerID:      peerId,
	}

	return &BlockChain{Header: &header, Blocks: make([]*Block, 0)}
}

//func (b BlockChain) GetLastBlock() *Block {
//	return b.LastBlock
//}
//
//func (b *BlockChain) AddBlock(blk *Block) (ret bool, err error){
//	if blk.Header.BlockStatus != Status_BLOCK_CONFIRMED {
//		ret = false
//		err = errors.New("unverified block")
//		return ret, err
//	}else if b.LastBlock != nil && b.LastBlock.Header.BlockHeight > 0 {
//		if b.LastBlock.Header.BlockHash != blk.Header.PreviousHash {
//			ret = false
//			err = errors.New("the hash values ​​of block and last Block are different")
//			return ret, err
//		}
//	}
//
//	BlkVarification, _ := blk.VerifyBlock()
//	// 모든 노드들에게 블록 검사를 보내고 결과를 받는 코드 작성해야 함.
//	//
//	if BlkVarification == false{ ret = false; err = errors.New("invalid block"); return ret, err }
//
//	b.Blocks = append(b.Blocks, blk)
//	b.LastBlock = blk
//	b.BlockNum++
//
//	return true, nil
//}
//
//func (b BlockChain) FindBlockIdx(hash string) (idx int, err error){
//	for idx = 0; idx < b.BlockNum; idx++{
//		if hash == b.Blocks[idx].Header.BlockHash{
//			return idx, nil
//		}
//	}
//	return -1, errors.New("BlockHash is not here")
//}
//
//func (b BlockChain) BlockLookUp(arg interface{}) (*Block, error) {
//	switch arg.(type)  {
//	case int:
//		idx := arg.(int)
//		return b.Blocks[idx], nil
//	case string:
//		idx, err := b.FindBlockIdx(arg.(string))
//		if err != nil { return Block{}, err }
//		return b.Blocks[idx], nil
//	default:
//		return nil, errors.New("no matched block")
//	}
//}

package blockchain

import (
	"it-chain/db/blockchaindb/blockchainleveldb"
	"errors"
)


type Ledger struct {
	LastBlock *Block
	Blocks    []*Block
	PeerId    string
	BlockNum  int
	DB        blockchainleveldb.BlockchainLevelDB
}

func (l *Ledger) CreateBlock(txList []*Transaction, createPeerId string) (*Block, error) {
	blk := CreateNewBlock(l.LastBlock, createPeerId)
	for _, tx := range txList {
		err := blk.PutTranscation(tx)
		if err != nil{ return nil, err }
	}
	blk.MakeMerkleTree()
	err := blk.GenerateBlockHash()
	if err != nil{ return nil, err }
	return blk, nil
}

func VerifyBlock(blk *Block) (bool, error){
	return blk.VerifyBlock()
}

func (l *Ledger) GetLastBlock() *Block{
	return l.LastBlock
}

func (l *Ledger) FindBlockIdx(hash string) (idx int, err error){
	for idx = 0; idx < l.BlockNum; idx++{
		if hash == l.Blocks[idx].Header.BlockHash{
			return idx, nil
		}
	}
	return -1, errors.New("BlockHash is not here")
}

func (l *Ledger) LookUpBlock(arg interface{}) (*Block, error) {
	switch arg.(type)  {
	case int:
		idx := arg.(int)
		return l.Blocks[idx], nil
	case string:
		idx, err := l.FindBlockIdx(arg.(string))
		if err != nil { return nil, err }
		return l.Blocks[idx], nil
	default:
		return nil, errors.New("no matched block")
	}
}

func (l *Ledger) AddBlock(blk *Block) (ret bool, err error){
	if blk.Header.BlockStatus != Status_BLOCK_CONFIRMED {
		ret = false
		err = errors.New("unverified block")
		return ret, err
	}else if l.LastBlock != nil && l.LastBlock.Header.BlockHeight > 0 {
		if l.LastBlock.Header.BlockHash != blk.Header.PreviousHash {
			ret = false
			err = errors.New("the hash values ​​of block and last Block are different")
			return ret, err
		}
	}

	BlkVarification, _ := blk.VerifyBlock()
	// 모든 노드들에게 블록 검사를 보내고 결과를 받는 코드 작성해야 함.
	//
	if BlkVarification == false{ ret = false; err = errors.New("invalid block"); return ret, err }

	l.Blocks = append(l.Blocks, blk)
	l.LastBlock = blk
	l.BlockNum++

	return true, nil
}
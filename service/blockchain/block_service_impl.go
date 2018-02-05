package blockchain

import (
	"it-chain/db/blockchaindb"
	"errors"
)

type Ledger struct {
	DB *blockchaindb.BlockchainDBImpl
}

func NewLedger(path string) BlockService{
	return &Ledger{DB:blockchaindb.CreateNewBlockchainDB(path)}
}

func (l *Ledger) CreateBlock(txList []*Transaction, createPeerId string) (*Block, error) {
	lastBlock, err := l.GetLastBlock()
	if err != nil { return nil, err }
	blk := CreateNewBlock(lastBlock, createPeerId)
	for _, tx := range txList {
		err = blk.PutTranscation(tx)
		if err != nil{ return nil, err }
	}
	blk.MakeMerkleTree()
	err = blk.GenerateBlockHash()
	if err != nil{ return nil, err }
	return blk, nil
}

func (l *Ledger) VerifyBlock(blk *Block) (bool, error){
	return blk.VerifyBlock()
}

func (l *Ledger) GetLastBlock() (*Block, error) {
	blk, err := l.DB.GetLastBlock()
	if err != nil{ return nil, err }
	return blk, err
}

func (l *Ledger) LookUpBlock(arg interface{}) (*Block, error) {
	switch arg.(type)  {
	case uint64:
		return l.DB.GetBlockByNumber(arg.(uint64))
	case string:
		return l.DB.GetBlockByHash(arg.(string))
	default:
		return nil, errors.New("no matched block")
	}
}

func (l *Ledger) AddBlock(blk *Block) (bool, error) {
	lastBlock, err := l.GetLastBlock()
	if err != nil { return false, err }
	if blk.Header.BlockStatus != Status_BLOCK_CONFIRMED {
		return false, errors.New("unverified block")
	}else if lastBlock != nil && lastBlock.Header.BlockHeight > 0 {
		if lastBlock.Header.BlockHash != blk.Header.PreviousHash {
			return false, errors.New("the hash values ​​of block and last Block are different")
		}
	}

	BlkVarification, _ := blk.VerifyBlock()

	if BlkVarification == false{ return false, errors.New("invalid block") }

	l.DB.AddBlock(blk)

	return true, nil
}
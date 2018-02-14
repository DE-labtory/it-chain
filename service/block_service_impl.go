package service

import (
	"it-chain/db/blockchaindb"
	"it-chain/domain"
	"errors"
)

type Ledger struct {
	DB *blockchaindb.BlockChainDB
}

func NewLedger(path string) BlockService{
	return &Ledger{DB:blockchaindb.CreateNewBlockchainDB(path)}
}

func (l *Ledger) CreateBlock(txList []*domain.Transaction, createPeerId string) (*domain.Block, error) {

	lastBlock, err := l.GetLastBlock()
	if err != nil { return nil, err }
	blk := domain.CreateNewBlock(lastBlock, createPeerId)
	for _, tx := range txList {
		err = blk.PutTranscation(tx)
		if err != nil{ return nil, err }
	}
	blk.MakeMerkleTree()
	err = blk.GenerateBlockHash()
	if err != nil{ return nil, err }
	return blk, nil
}

func (l *Ledger) VerifyBlock(blk *domain.Block) (bool, error){

	lastBlock, err := l.DB.GetLastBlock()

	if err != nil {
		return false, err
	}

	if lastBlock.Header.BlockHeight + 1 != blk.Header.BlockHeight{
		return false, errors.New("Block height misMatched")
	}

	if lastBlock.Header.BlockHash != blk.Header.PreviousHash{
		return false, errors.New("Block hash is different")
	}

	_, err = blk.VerifyBlock()

	if err != nil {
		return false, err
	}

	blk.Header.BlockStatus = domain.Status_BLOCK_CONFIRMED
	return true, nil
}

func (l *Ledger) GetLastBlock() (*domain.Block, error) {
	blk, err := l.DB.GetLastBlock()
	if err != nil{ return nil, err }
	return blk, err
}

func (l *Ledger) LookUpBlock(arg interface{}) (*domain.Block, error) {
	switch arg.(type)  {
	case int:
		arg = uint64(arg.(int))
		return l.DB.GetBlockByNumber(arg.(uint64))
	case string:
		return l.DB.GetBlockByHash(arg.(string))
	default:
		return nil, errors.New("no matched block")
	}
}

func (l *Ledger) AddBlock(blk *domain.Block) (bool, error) {
	lastBlock, err := l.GetLastBlock()
	if err != nil { return false, err }
	if blk.Header.BlockStatus != domain.Status_BLOCK_CONFIRMED {
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
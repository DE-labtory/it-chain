package service

import (
	"github.com/it-chain/it-chain-Engine/db/blockchaindb"
	"github.com/it-chain/it-chain-Engine/domain"
	"errors"
)

type Ledger struct {
	DB blockchaindb.BlockChainDB
}

func NewLedger(path string) BlockService{
	return &Ledger{DB:blockchaindb.CreateNewBlockchainDB(path)}
}

func (l *Ledger) Close() {
	l.DB.Close()
}

func (l *Ledger) CreateBlock(txList []*domain.Transaction, createPeerId string) (*domain.Block, error) {

	lastBlock, err := l.GetLastBlock()

	if err != nil {
		return nil, err
	}

	if lastBlock == nil {
		lastBlock, err = l.CreateGenesisBlock()
	}

	blk := domain.CreateNewBlock(lastBlock, createPeerId)

	for _, tx := range txList {
		err = blk.PutTranscation(tx)
		if err != nil{
			return nil, err
		}
	}

	blk.MakeMerkleTree()
	err = blk.GenerateBlockHash()

	if err != nil{
		return nil, err
	}

	return blk, nil
}

func (l *Ledger) VerifyBlock(blk *domain.Block) (bool, error){

	lastBlock, err := l.DB.GetLastBlock()

	if lastBlock == nil {
		_, err = blk.VerifyBlock()

		if err != nil {
			return false, err
		}

		blk.Header.BlockStatus = domain.Status_BLOCK_CONFIRMED
		return true, nil
	}

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

	if err != nil{
		return nil, err
	}

	if blk == nil{

		return l.CreateGenesisBlock()
	}

	return blk, err
}

func (l *Ledger) LookUpBlock(arg interface{}) (*domain.Block, error) {
	_, err := l.GetLastBlock()
	if err != nil{
		return nil, errors.New("no block exist")
	}
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
	_, err := l.LookUpBlock(blk.Header.BlockHash)
	if err == nil{ return false, errors.New("already this blk exist") }

	lastBlock, err := l.GetLastBlock()

	if err != nil {
		return false, err
	}
	if blk.Header.BlockStatus != domain.Status_BLOCK_CONFIRMED {
		return false, errors.New("unverified block")
	}else if lastBlock != nil && lastBlock.Header.BlockHeight > 0 {

		if lastBlock.Header.BlockHash != blk.Header.PreviousHash {
			return false, errors.New("the hash values ​​of block and last Block are different")
		}
	}

	BlkVarification, _ := blk.VerifyBlock()

	if BlkVarification == false{
		return false, errors.New("invalid block")
	}

	l.DB.AddBlock(blk)

	return true, nil
}

func (l *Ledger) CreateGenesisBlock() (*domain.Block, error) {

	genesisBlock := domain.CreateNewBlock(nil, "0")
	genesisBlock.MakeMerkleTree()
	genesisBlock.GenerateBlockHash()
	BlkVarification, _ := genesisBlock.VerifyBlock()

	if BlkVarification == false{
		return nil, errors.New("invalid block")
	}

	l.DB.AddBlock(genesisBlock)
	return genesisBlock, nil
}
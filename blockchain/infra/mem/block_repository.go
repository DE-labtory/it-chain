package mem

import (
	"sync"

	"errors"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/yggdrasill"
)

var ErrAddBlock = errors.New("Error in adding block")
var ErrGetBlock = errors.New("Error in getting block")
var ErrEmptyBlock = errors.New("Error when block is empty that should be not")
var ErrNewBlockStorage = errors.New("Error in constructing block storage")

type blockRepository struct {
	mux *sync.RWMutex
	yggdrasill.BlockStorageManager
}

func NewBlockRepositoryImpl(dbPath string) (*blockRepository, error) {
	validator := new(blockchain.DefaultValidator)
	db := leveldbwrapper.CreateNewDB(dbPath)
	opts := map[string]interface{}{}

	blockStorage, err := yggdrasill.NewBlockStorage(db, validator, opts)
	if err != nil {
		return nil, ErrNewBlockStorage
	}

	return &blockRepository{
		mux:                 &sync.RWMutex{},
		BlockStorageManager: blockStorage,
	}, nil
}

func (br *blockRepository) Save(block blockchain.DefaultBlock) error {
	br.mux.Lock()
	defer br.mux.Unlock()

	err := br.BlockStorageManager.AddBlock(&block)
	if err != nil {
		return ErrAddBlock
	}

	return nil
}

func (br *blockRepository) FindLast() (blockchain.DefaultBlock, error) {
	br.mux.Lock()
	defer br.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := br.BlockStorageManager.GetLastBlock(block)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetBlock
	}

	return *block, nil
}
func (br *blockRepository) FindByHeight(height uint64) (blockchain.DefaultBlock, error) {
	br.mux.Lock()
	defer br.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := br.BlockStorageManager.GetBlockByHeight(block, height)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetBlock
	}

	return *block, nil
}

func (br *blockRepository) FindAll() ([]blockchain.DefaultBlock, error) {
	br.mux.Lock()
	defer br.mux.Unlock()

	blocks := []blockchain.DefaultBlock{}

	// set
	lastBlock := &blockchain.DefaultBlock{}

	err := br.BlockStorageManager.GetLastBlock(lastBlock)

	if err != nil {
		return nil, err
	}

	// check empty
	if lastBlock.IsEmpty() {
		return blocks, nil
	}

	lastHeight := lastBlock.GetHeight()

	// get blocks
	for i := uint64(0); i <= lastHeight; i++ {

		block := &blockchain.DefaultBlock{}

		err := br.BlockStorageManager.GetBlockByHeight(block, i)

		if err != nil {
			return nil, err
		}

		if block.IsEmpty() {
			return nil, ErrEmptyBlock
		}

		blocks = append(blocks, *block)
	}

	return blocks, nil
}

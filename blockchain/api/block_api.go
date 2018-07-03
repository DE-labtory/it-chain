package api

import (
	"encoding/json"
	"time"

	"io/ioutil"
	"os"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"fmt"
	"github.com/pkg/errors"
)

type BlockRepository interface {
	GetValidator() blockchain.Validator
	GetLastBlock(block blockchain.Block) error
	AddBlock(block blockchain.Block) error
	NewEmptyBlock() (blockchain.Block, error)
	Close()
}

var ErrNilBlock = errors.New("block is nil")

type BlockApi struct {
	blockRepository BlockRepository
	publisherId          string
	blockPool blockchain.BlockPool
}

func NewBlockApi(blockRepository BlockRepository, publisherId string, blockPool blockchain.BlockPool) (BlockApi, error) {
	return BlockApi{
		blockRepository: blockRepository,
		publisherId:          publisherId,
		blockPool: blockPool,
	}, nil
}

// TODO: 테스트 필요.
func (bApi *BlockApi) CreateGenesisBlock(genesisConfFilePath string) (blockchain.Block, error) {
	byteValue, err := getConfigFromJson(genesisConfFilePath)
	if err != nil {
		return nil, err
	}

	validator := bApi.blockRepository.GetValidator()

	var GenesisBlock blockchain.Block

	json.Unmarshal(byteValue, &GenesisBlock)
	GenesisBlock.SetTimestamp((time.Now()).Round(0))
	Seal, err := validator.BuildSeal(GenesisBlock)
	if err != nil {
		return nil, err
	}

	GenesisBlock.SetSeal(Seal)
	return GenesisBlock, nil
}

func (bApi *BlockApi) CreateBlock(txList []blockchain.Transaction) (blockchain.Block, error) {
	repo := bApi.blockRepository

	block, err := repo.NewEmptyBlock()
	if err != nil {
		return nil, err
	}

	v := bApi.blockRepository.GetValidator()

	txSeal, err := v.BuildTxSeal(txList)
	if err != nil {
		return nil, err
	}

	block.SetTxSeal(txSeal)

	for _, tx := range txList {
		block.PutTx(tx)
	}

	block.SetTimestamp(time.Now())

	blockSeal, err := v.BuildSeal(block)

	block.SetSeal(blockSeal)

	return block, nil
}

func getConfigFromJson(filePath string) ([]uint8, error) {
	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}

// TODO: Check 과정에서 임의의 노드에게서 받은 blockchain 정보로 동기화 되었는지 확인한다.
func (bApi *BlockApi) SyncedCheck(block blockchain.Block) error {
	return nil
}

// 받은 block을 block pool에 추가한다.
func (bApi *BlockApi) AddBlockToPool(block blockchain.Block) {
	if block == nil {
		fmt.Println("block is nil")
		return
	}
	bApi.blockPool.Add(block)
}

// TODO
func (bApi *BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	// Get block from pool
	block := bApi.blockPool.Get(height)
	if block == nil {
		return ErrNilBlock
	}

	// Get my last block
	lastBlock := &blockchain.DefaultBlock{}
	bApi.blockRepository.GetLastBlock(lastBlock)

	// Compare height
	if block.GetHeight() > lastBlock.GetHeight() + 1 {
		// TODO: Start synchronize

	} else if block.GetHeight() == lastBlock.GetHeight() + 1 {
		// Save
		bApi.blockRepository.AddBlock(block)

		bApi.blockPool.Delete(height)

	} else {
		// Got shorter height block, but this is not an error
		fmt.Printf("got shorter height block [%d < %d]", block.GetHeight(), lastBlock.GetHeight());
	}

	return nil
}
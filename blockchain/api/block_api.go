package api

import (
	"encoding/json"
	"time"

	"io/ioutil"
	"os"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
)

type BlockApi struct {
	blockchainRepository blockchain.Repository
	eventRepository      *midgard.Repository
	publisherId          string
}

func NewBlockApi(blockchainRepository blockchain.Repository, eventRepository *midgard.Repository, publisherId string) (BlockApi, error) {
	return BlockApi{
		blockchainRepository: blockchainRepository,
		eventRepository:      eventRepository,
		publisherId:          publisherId,
	}, nil
}

// TODO: 테스트 필요.
func (bApi *BlockApi) CreateGenesisBlock(genesisConfFilePath string) (blockchain.Block, error) {
	byteValue, err := getConfigFromJson(genesisConfFilePath)
	if err != nil {
		return nil, err
	}

	validator := bApi.blockchainRepository.GetValidator()

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
	repo := bApi.blockchainRepository

	block, err := repo.NewEmptyBlock()
	if err != nil {
		return nil, err
	}

	v := bApi.blockchainRepository.GetValidator()

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

package api

import (
	"encoding/json"
	"time"

	"io/ioutil"
	"os"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
	"github.com/it-chain/yggdrasill/common"
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

func (bApi *BlockApi) CreateGenesisBlock(genesisConfFilePath string) (common.Block, error) {
	byteValue, err := getConfigFromJson(genesisConfFilePath)
	if err != nil {
		return nil, err
	}

	validator := bApi.blockchainRepository.GetValidator()

	var GenesisBlock common.Block

	json.Unmarshal(byteValue, &GenesisBlock)
	GenesisBlock.SetTimestamp((time.Now()).Round(0))
	Seal, err := validator.BuildSeal(GenesisBlock)
	if err != nil {
		return nil, err
	}

	GenesisBlock.SetSeal(Seal)
	return GenesisBlock, nil
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

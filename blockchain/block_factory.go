package blockchain

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

var ErrGetConfig = errors.New("error when get Config")
var ErrBuildingSeal = errors.New("error when building Seal")
var ErrBuildingTxSeal = errors.New("error when building TxSeal")

func CreateGenesisBlock(genesisconfFilePath string) (Block, error) {

	GenesisBlock := &DefaultBlock{}
	validator := DefaultValidator{}

	byteValue, err := configFromJson(genesisconfFilePath)

	if err != nil {
		return nil, ErrGetConfig
	}

	json.Unmarshal(byteValue, &GenesisBlock)

	GenesisBlock.SetPrevSeal(GenesisBlock.PrevSeal)
	GenesisBlock.SetHeight(GenesisBlock.Height)
	GenesisBlock.SetTxSeal(GenesisBlock.TxSeal)
	GenesisBlock.SetCreator(GenesisBlock.Creator)
	GenesisBlock.SetTimestamp((time.Now()).Round(0))

	Seal, err := validator.BuildSeal(GenesisBlock)

	if err != nil {
		return nil, ErrBuildingSeal
	}

	GenesisBlock.SetSeal(Seal)

	return GenesisBlock, nil
}

func CreateProposedBlock(prevSeal []byte, height uint64, txList []Transaction, Creator []byte) (Block, error) {

	Block := &DefaultBlock{}
	validator := DefaultValidator{}

	Block.SetPrevSeal(prevSeal)
	Block.SetHeight(height)
	Block.SetCreator(Creator)

	for _, tx := range txList {
		Block.PutTx(tx)
	}

	txSeal, err := validator.BuildTxSeal(Block.GetTxList())

	if err != nil {
		return nil, ErrBuildingTxSeal
	}

	Block.SetTxSeal(txSeal)

	Block.SetTimestamp((time.Now()).Round(0))

	Seal, err := validator.BuildSeal(Block)

	if err != nil {
		return nil, ErrBuildingSeal
	}

	Block.SetSeal(Seal)

	return Block, nil
}

func configFromJson(filePath string) ([]uint8, error) {
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

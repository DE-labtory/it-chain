package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

func CreateGenesisBlock(genesisconfFilePath string) (*DefaultBlock, error) {
	byteValue, err := ConfigFromJson(genesisconfFilePath)
	if err != nil {
		return nil, err
	}

	validator := new(DefaultValidator)

	var GenesisBlock *DefaultBlock

	json.Unmarshal(byteValue, &GenesisBlock)
	GenesisBlock.SetTimestamp((time.Now()).Round(0))
	Seal, err := validator.BuildSeal(GenesisBlock)
	if err != nil {
		return nil, err
	}

	GenesisBlock.SetSeal(Seal)
	GenesisBlock.SetPrevSeal(GenesisBlock.PrevSeal)
	GenesisBlock.SetHeight(GenesisBlock.Height)
	GenesisBlock.SetTxSeal(GenesisBlock.TxSeal)
	GenesisBlock.SetCreator(GenesisBlock.Creator)
	return GenesisBlock, nil
}

func CreateBlock(txList []Transaction) (Block, error) {

	var Block *DefaultBlock

	validator := new(DefaultValidator)

	txSeal, err := validator.BuildTxSeal(txList)

	if err != nil {
		return nil, err
	}

	Block.SetTxSeal(txSeal)

	for _, tx := range txList {
		Block.PutTx(tx)
	}

	Seal, err := validator.BuildSeal(Block)
	Block.SetSeal(Seal)
	Block.SetTimestamp((time.Now()).Round(0))

	return Block, nil
}

func ConfigFromJson(filePath string) ([]uint8, error) {
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

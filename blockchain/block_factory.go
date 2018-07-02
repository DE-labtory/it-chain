package blockchain

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

var ErrGetConfig = errors.New("error when get Config")

func CreateGenesisBlock(genesisconfFilePath string) (*DefaultBlock, error) {

	byteValue, err := configFromJson(genesisconfFilePath)

	if err != nil {
		return nil, ErrGetConfig
	}

	validator := new(DefaultValidator)

	var GenesisBlock *DefaultBlock

	json.Unmarshal(byteValue, &GenesisBlock)

	GenesisBlock.SetPrevSeal(GenesisBlock.PrevSeal)
	GenesisBlock.SetHeight(GenesisBlock.Height)
	GenesisBlock.SetTxSeal(GenesisBlock.TxSeal)
	GenesisBlock.SetCreator(GenesisBlock.Creator)
	GenesisBlock.SetTimestamp((time.Now()).Round(0))

	Seal, err := validator.BuildSeal(GenesisBlock)

	if err != nil {
		return nil, err
	}

	GenesisBlock.SetSeal(Seal)

	return GenesisBlock, nil
}

//func CreateBlock(prevSeal []byte, Height uint64, txList []Transaction, Timestamp time.Time, Creator []byte) (Block, error) {
//
//	var Block *DefaultBlock
//
//	validator := new(DefaultValidator)
//
//	txSeal, err := validator.BuildTxSeal(txList)
//
//	if err != nil {
//		return nil, err
//	}
//
//	Block.SetTxSeal(txSeal)
//
//	for _, tx := range txList {
//		Block.PutTx(tx)
//	}
//
//	Seal, err := validator.BuildSeal(Block)
//	Block.SetSeal(Seal)
//	Block.SetTimestamp((time.Now()).Round(0))
//
//	return Block, nil
//}

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

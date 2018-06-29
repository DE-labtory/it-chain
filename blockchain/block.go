package blockchain

import (
	"encoding/json"
	"time"

	"io/ioutil"
	"os"

	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/impl"
)

type Block = common.Block

type DefaultBlock = impl.DefaultBlock

type BlockRepository interface {
	yggdrasill.BlockStorageManager
	NewEmptyBlock() (Block, error)
	GetBlockCreator() string
}

func CreateGenesisBlock(genesisconfFilePath string) (*impl.DefaultBlock, error) {
	byteValue, err := ConfigFromJson(genesisconfFilePath)
	if err != nil {
		return nil, err
	}

	validator := new(impl.DefaultValidator)

	var GenesisBlock *impl.DefaultBlock

	json.Unmarshal(byteValue, &GenesisBlock)
	GenesisBlock.SetTimestamp((time.Now()).Round(0))
	Seal, err := validator.BuildSeal(GenesisBlock)
	if err != nil {
		return nil, err
	}

	GenesisBlock.SetSeal(Seal)
	GenesisBlock.SetPrevSeal(GenesisBlock.PrevSeal)
	GenesisBlock.SetHeight(GenesisBlock.Height)
	GenesisBlock.SetHeight(GenesisBlock.Height)
	GenesisBlock.SetTxSeal(GenesisBlock.TxSeal)
	GenesisBlock.SetCreator(GenesisBlock.Creator)
	return GenesisBlock, nil
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

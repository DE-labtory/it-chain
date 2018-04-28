package service

import (
	"encoding/json"
	"time"

	"github.com/it-chain/yggdrasill/block"
)

func CreateGenesisBlock(genesisconfFilePath string) (*block.DefaultBlock, error) {
	byteValue, err := ConfigFromJson(genesisconfFilePath)
	if err != nil {
		return nil, err
	}

	var GenesisBlock *block.DefaultBlock
	json.Unmarshal(byteValue, &GenesisBlock)
	GenesisBlock.Header.TimeStamp = time.Now()
	return GenesisBlock, nil
}

package service

import (
	"encoding/json"
	"time"

	"github.com/it-chain/yggdrasill/impl"
)

func CreateGenesisBlock(genesisconfFilePath string) (*impl.DefaultBlock, error) {
	byteValue, err := ConfigFromJson(genesisconfFilePath)
	if err != nil {
		return nil, err
	}
	var GenesisBlock *impl.DefaultBlock
	json.Unmarshal(byteValue, &GenesisBlock)
	GenesisBlock.Timestamp = time.Now()
	return GenesisBlock, nil
}

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
	//ToDo: yggdrasill의 Update사항을 반영하여, validate과정을 추가하여야 합니다.(주석 작성자 GitID:junk-sound)

	var GenesisBlock *impl.DefaultBlock
	json.Unmarshal(byteValue, &GenesisBlock)
	GenesisBlock.Timestamp = time.Now()
	return GenesisBlock, nil
}

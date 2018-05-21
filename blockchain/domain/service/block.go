package service

import (
	"encoding/json"

	"time"

	"fmt"

	"github.com/it-chain/yggdrasill/impl"
)

func CreateGenesisBlock(genesisconfFilePath string) (*impl.DefaultBlock, error) {
	fmt.Println("Service/block.go")
	byteValue, err := ConfigFromJson(genesisconfFilePath)
	if err != nil {
		return nil, err
	}

	validator := new(impl.DefaultValidator)

	var GenesisBlock *impl.DefaultBlock

	json.Unmarshal(byteValue, &GenesisBlock)
	GenesisBlock.SetTimestamp((time.Now()))
	Seal, err := validator.BuildSeal(GenesisBlock)
	if err != nil {
		return nil, err
	}

	GenesisBlock.SetSeal(Seal)
	fmt.Println("Seal:")
	fmt.Println(Seal)

	GenesisBlock.SetPrevSeal(GenesisBlock.PrevSeal)
	GenesisBlock.SetHeight(GenesisBlock.Height)
	GenesisBlock.SetHeight(GenesisBlock.Height)
	GenesisBlock.SetTxSeal(GenesisBlock.TxSeal)
	GenesisBlock.SetCreator(GenesisBlock.Creator)
	return GenesisBlock, nil
}

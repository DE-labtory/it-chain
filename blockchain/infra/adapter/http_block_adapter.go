package adapter

import (
	"encoding/json"

	"github.com/go-resty/resty"
	"github.com/it-chain/engine/blockchain"
)

type HttpBlockAdapter struct {
}

func (a HttpBlockAdapter) GetLastBlock(address string) (blockchain.DefaultBlock, error) {

	resp, err := resty.R().
		SetQueryString("last=true").
		SetHeader("Accept", "application/json").
		SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get("http://" + address + "/blocks")

	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	block := blockchain.DefaultBlock{}

	if err := json.Unmarshal(resp.Body(), &block); err != nil {
		return blockchain.DefaultBlock{}, err
	}

	return block, nil
}

func (a HttpBlockAdapter) GetBlockByHeight(address string, height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	resp, err := resty.R().Get("http://" + address + "blocks?height=" + string(height))
	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	block := blockchain.DefaultBlock{}

	if err := json.Unmarshal(resp.Body(), block); err != nil {
		return blockchain.DefaultBlock{}, err
	}

	return block, nil
}

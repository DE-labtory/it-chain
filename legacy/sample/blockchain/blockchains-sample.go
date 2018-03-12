package main

import (
	"github.com/it-chain/it-chain-Engine/legacy/service/blockchain"
	"fmt"
)

func main(){
	var block = blockchain.CreateNewBlockChain("channel1","peer0")
	block.Blocks = append(block.Blocks, &blockchain.Block{})
	fmt.Print(len(block.Blocks))
}
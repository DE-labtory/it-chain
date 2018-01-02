package main

import (
	"it-chain/service"
	"fmt"
)

func main(){
	var block = service.CreateNewBlockChain("channel1","peer0")
	block.Blocks = append(block.Blocks, &service.Block{})
	fmt.Print(len(block.Blocks))
}
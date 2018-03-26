package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"strconv"
	"fmt"
	"github.com/it-chain/it-chain-Engine/common"
)

const txsize = 999

func TestCreateNewBlockTest(t *testing.T){
	var block = CreateNewBlock(nil)
	assert.Equal(t, 0, block.BlockData.TransactionCount)
}


func TestBlock_MakeMerkleTree(t *testing.T) {
	block := CreateNewBlock(nil)
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction("1", strconv.Itoa(i), time.Now(), SetTxData("", Invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}
	block.MakeMerkleTree()

	for h := block.BlockData.MerkleTreeHeight-1; h >= 0; h--{
		for _, that := range block.BlockData.MerkleTree[h]{
			if that != ""{
				//fmt.Printf("0 ")
				fmt.Printf("%s ", that)
			}
		}
		fmt.Printf("\n")
	}

	assert.Equal(t, 999, block.BlockData.TransactionCount, "tx_count")

	str := []string{block.BlockData.MerkleTree[0][len(block.BlockData.MerkleTree[0])-1], block.BlockData.MerkleTree[0][len(block.BlockData.MerkleTree[0])-1]}
	assert.Equal(t, block.BlockData.MerkleTree[1][len(block.BlockData.MerkleTree[1])-1], common.ComputeSHA256(str), "hash")


	//assert.Equal(t, 16, len(block.MerkleTree[0]), "tx_list")
	assert.Equal(t, 11, block.BlockData.MerkleTreeHeight, "mt_height")
}

func TestBlock_MakeMerklePath(t *testing.T) {
	block := CreateNewBlock(nil)
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction("1", strconv.Itoa(i), time.Now(), SetTxData("", Invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}
	block.MakeMerkleTree()

	idx := 2

	path := block.MakeMerklePath(idx)
	hash := block.BlockData.Transactions[idx].TransactionHash
	for _, sibling_hash := range path{
		str := []string{hash, sibling_hash}
		hash = common.ComputeSHA256(str)
	}

	assert.Equal(t, block.BlockHeader.MerkleTreeRootHash, hash)
}

func TestBlock_GenerateBlockHash(t *testing.T) {
	block := CreateNewBlock(nil)
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction("1", strconv.Itoa(i), time.Now(), SetTxData("", Invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}
	block.MakeMerkleTree()

	assert.NoError(t, block.GenerateBlockHash())
}

func TestSerializationAndDeserialization(t *testing.T) {
	block := CreateNewBlock(nil)
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction("1", strconv.Itoa(i), time.Now(), SetTxData("", Invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}
	block.MakeMerkleTree()

	str, err1 := block.BlockSerialize()
	assert.NoError(t, err1)

	_, err2 := BlockDeserialize(str)

	assert.NoError(t, err2)
}

func TestVerifyBlock(t *testing.T){
	block := CreateNewBlock(nil)
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction("1", strconv.Itoa(i), time.Now(), SetTxData("", Invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}

	block.MakeMerkleTree()

	_, err := block.VerifyBlock()
	assert.NoError(t, err)
}
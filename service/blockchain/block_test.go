package blockchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"strconv"
	"fmt"
	"it-chain/common"
)

const txsize = 999

func TestCreateNewBlockTest(t *testing.T){

	var block = CreateNewBlock(nil, "")
	assert.Equal(t, 0, block.TransactionCount)
}

func TestBlock_PutTranscation(t *testing.T) {
	block := CreateNewBlock(nil, "")
	tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), general, time.Now(), SetTxData("", invoke, SetTxMethodParameters(0, "", []string{""}), ""))
	tx.GenerateHash()
	err := block.PutTranscation(tx)
	assert.NoError(t, err)
}

func TestBlock_MakeMerkleTree(t *testing.T) {
	block := CreateNewBlock(nil, "")
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), general, time.Now(), SetTxData("", invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}
	block.MakeMerkleTree()

	for h := block.MerkleTreeHeight-1; h >= 0; h--{
		for _, that := range block.MerkleTree[h]{
			if that != ""{
				//fmt.Printf("0 ")
				fmt.Printf("%s ", that)
			}
		}
		fmt.Printf("\n")
	}

	assert.Equal(t, 999, block.TransactionCount, "tx_count")

	str := []string{block.MerkleTree[0][len(block.MerkleTree[0])-1], block.MerkleTree[0][len(block.MerkleTree[0])-1]}
	assert.Equal(t, block.MerkleTree[1][len(block.MerkleTree[1])-1], common.ComputeSHA256(str), "hash")


	//assert.Equal(t, 16, len(block.MerkleTree[0]), "tx_list")
	assert.Equal(t, 11, block.MerkleTreeHeight, "mt_height")
}

func TestBlock_FindTransactionIndex(t *testing.T) {
	block := CreateNewBlock(nil, "")
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), general, time.Now(), SetTxData("", invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}
	block.MakeMerkleTree()

	idx, err := block.FindTransactionIndex(block.MerkleTree[0][len(block.MerkleTree[0])-1])

	assert.NoError(t, err)
	assert.Equal(t, 998, idx)

}

func TestBlock_MakeMerklePath(t *testing.T) {
	block := CreateNewBlock(nil, "")
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), general, time.Now(), SetTxData("", invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}
	block.MakeMerkleTree()

	idx := 2

	path := block.MakeMerklePath(idx)
	hash := block.Transactions[idx].TransactionHash
	for _, sibling_hash := range path{
		str := []string{hash, sibling_hash}
		hash = common.ComputeSHA256(str)
	}

	assert.Equal(t, block.Header.MerkleTreeRootHash, hash)
}

func TestBlock_GenerateBlockHash(t *testing.T) {
	block := CreateNewBlock(nil, "")
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), general, time.Now(), SetTxData("", invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}
	block.MakeMerkleTree()

	assert.NoError(t, block.GenerateBlockHash())
}

func TestSerializationAndDeserialization(t *testing.T) {
	block := CreateNewBlock(nil, "12312313")
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), general, time.Now(), SetTxData("", invoke, SetTxMethodParameters(0, "", []string{""}), ""))
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
	block := CreateNewBlock(nil, "12312313")
	for i := 0; i < txsize; i++{
		tx := CreateNewTransaction(strconv.Itoa(1), strconv.Itoa(1), general, time.Now(), SetTxData("", invoke, SetTxMethodParameters(0, "", []string{""}), ""))
		tx.GenerateHash()
		block.PutTranscation(tx)
	}

	block.MakeMerkleTree()

	_, err := block.VerifyBlock()
	assert.NoError(t, err)
}
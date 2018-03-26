package model

import (
	"time"
	"github.com/it-chain/it-chain-Engine/common"
	"errors"
	"strconv"
)

type Status_BLOCK int

const (
	Status_BLOCK_UNCONFIRMED Status_BLOCK = iota
	Status_BLOCK_CONFIRMED
)

type Block struct {
	BlockHeader *BlockHeader
	BlockData   *BlockData
}

type BlockData struct {
	MerkleTree       [][]string
	Transactions     []*Transaction
	TxIdx            map[string]int
	MerkleTreeHeight int
	TransactionCount int
}

type BlockHeader struct {
	Number             uint64
	PreviousHash       string
	Version            string
	MerkleTreeRootHash string
	TimeStamp          time.Time
	BlockHeight        int
	BlockStatus        Status_BLOCK
	BlockHash          string
}

func CreateNewBlock(prevBlock *Block) *Block{
	var header BlockHeader
	if prevBlock == nil{
		header.Number = 0
		header.PreviousHash = ""
		header.Version = ""
		header.BlockHeight = 0
	} else {
		header.Number = prevBlock.BlockHeader.Number + 1
		header.PreviousHash = prevBlock.BlockHeader.BlockHash
		header.Version = prevBlock.BlockHeader.Version
		header.BlockHeight = prevBlock.BlockHeader.BlockHeight + 1
	}
	header.TimeStamp = time.Now().Round(0)
	header.BlockStatus = Status_BLOCK_UNCONFIRMED

	var data BlockData

	data.MerkleTree = make([][]string, 0)
	data.MerkleTreeHeight = 0
	data.TransactionCount = 0
	data.Transactions = make([]*Transaction, 0)
	data.TxIdx = make(map[string]int)

	return &Block{BlockHeader:&header, BlockData:&data}
}

//func (block *Block) PutTranscation(txList []*Transaction) {
//	for _, tx := range txList{
//		block.BlockData.Transactions = append(block.BlockData.Transactions, tx)
//		block.BlockData.TxIdx[tx.TransactionHash] = block.BlockData.TransactionCount
//		block.BlockData.TransactionCount++
//	}
//}

func (block *Block) PutTranscation(input interface{}) {
	switch v := input.(type) {
	case *Transaction:
		block.BlockData.Transactions = append(block.BlockData.Transactions, v)
		block.BlockData.TxIdx[v.TransactionHash] = block.BlockData.TransactionCount
		block.BlockData.TransactionCount++
	case []*Transaction:
		for _, tx := range v{
			block.BlockData.Transactions = append(block.BlockData.Transactions, tx)
			block.BlockData.TxIdx[tx.TransactionHash] = block.BlockData.TransactionCount
			block.BlockData.TransactionCount++
		}
	}
}

func (block *Block) MakeMerkleTree(){
	if block.BlockData.TransactionCount == 0 {
		block.BlockHeader.MerkleTreeRootHash = ""
		return
	}

	var mtList []string

	for _, h := range block.BlockData.Transactions{
		mtList = append(mtList, h.TransactionHash)
	}
	for {
		treeLength := len(mtList)
		block.BlockData.MerkleTreeHeight++
		if treeLength <= 1 {
			block.BlockData.MerkleTree = append(block.BlockData.MerkleTree, mtList)
			break
		} else if treeLength % 2 == 1 {
			mtList = append(mtList, mtList[treeLength - 1])
			treeLength++
		}
		block.BlockData.MerkleTree = append(block.BlockData.MerkleTree, mtList)
		var tmpMtList []string
		for x := 0; x < treeLength/2; x++{
			idx := x * 2
			hashArg := []string{mtList[idx], mtList[idx+1]}
			mkHash := common.ComputeSHA256(hashArg)
			tmpMtList = append(tmpMtList, mkHash)
		}
		mtList = tmpMtList
	}
	if len(mtList) == 1 {
		block.BlockHeader.MerkleTreeRootHash = mtList[0]
	}
}

func (block Block) MakeMerklePath(idx int) (path []string){
	for i := 0; i < block.BlockData.MerkleTreeHeight-1; i++{
		path = append(path, block.BlockData.MerkleTree[i][(idx >> uint(i)) ^ 1])
	}
	return path
}

func (block *Block) GenerateBlockHash() error{
	if block.BlockHeader.MerkleTreeRootHash == "" {
		return errors.New("no merkle tree root hash")
	}

	str := []string{block.BlockHeader.MerkleTreeRootHash, block.BlockHeader.TimeStamp.String(), block.BlockHeader.PreviousHash}
	block.BlockHeader.BlockHash = common.ComputeSHA256(str)
	return nil
}

func (block Block) BlockSerialize() ([]byte, error){
	return common.Serialize(block)
}

func BlockDeserialize(by []byte) (Block, error) {
	block := Block{}
	err := common.Deserialize(by, &block)
	return block, err
}

// 해당 트랜잭션이 정당한지 머클패스로 검사함
func (block Block) VerifyTx(tx Transaction) (bool, error) {
	if block.BlockHeader.BlockHeight == 0 && block.BlockData.TransactionCount == 0 {
		return true, nil;
	}

	hash := tx.TransactionHash
	idx := block.BlockData.TxIdx[hash]
	merklePath := block.MakeMerklePath(idx)

	for _, sibling_hash := range merklePath{
		str := []string{hash, sibling_hash}
		hash = common.ComputeSHA256(str)
	}

	if hash == block.BlockHeader.MerkleTreeRootHash{
		return true, nil
	} else {
		return false, errors.New("tx is invalid")
	}
}

// 블럭내의 모든 트랜잭션들이 정당한지 머클패스로 검사함
func (block Block) VerifyBlock() (bool, error) {
	for idx := 0; idx < block.BlockData.TransactionCount; idx++{
		txVarification, txErr := block.VerifyTx(*block.BlockData.Transactions[idx])
		if txVarification == false  {
			err := errors.New("block is invalid --- " + strconv.Itoa(idx) + "'s " + txErr.Error())
			return false, err
		}
	}
	return true, nil
}
package blockchain

import (
	"time"
	"errors"
	"it-chain/common"
)

type Status int

const (
	Status_BLOCK_UNCONFIRMED Status = 0
	Status_BLOCK_CONFIRMED   Status = 1
)

type Block struct {
	Header          	*BlockHeader
	MerkleTree      	[][]string
	MerkleTreeHeight	int
	TransactionCount	int
	Transactions	 	[]*Transaction
}

func (s *Block) Reset() { *s = Block{} }

type BlockHeader struct {
	Number             uint64
	PreviousHash       string
	Version            string
	MerkleTreeRootHash string
	TimeStamp          time.Time
	BlockHeight        int
	BlockStatus        Status
	CreatedPeerID      string
	Signature          []byte
	PeerId             string
}

func CreateNewBlock(num uint64, prev_hash string, ver string, t time.Time, create_peerId string) *Block{
	var header = BlockHeader{
		Number: num,
		PreviousHash: prev_hash,
		Version:      ver,
		TimeStamp:	t,
		BlockHeight: 0,
		BlockStatus: Status_BLOCK_UNCONFIRMED,
		CreatedPeerID: create_peerId,
	}
	return &Block{Header:&header, MerkleTree:make([][]string, 0), TransactionCount:0, Transactions:make([]*Transaction, 0)}
}

func (s *BlockHeader) Reset() { *s = BlockHeader{} }

func (s *Block) PutTranscation(tx *Transaction) bool{
	if tx.TransactionStatus == Status_BLOCK_UNCONFIRMED{
		if tx.Validate(){
			tx.TransactionStatus = Status_BLOCK_CONFIRMED
		} else {
			return false
		}
	}
	s.Transactions = append(s.Transactions, tx)
	s.TransactionCount++
	return true
}

func (s Block) FindTransactionIndex(hash string) (idx int, err error){
	for idx = 0; idx < s.TransactionCount; idx++{
		if hash == s.Transactions[idx].TransactionHash{
			return idx, nil
		}
	}
	return -1, errors.New("wrong Hash")
}

func (s *Block) MakeMerkleTree(){
	var mtList []string
	for _, h := range s.Transactions{
		mtList = append(mtList, h.TransactionHash)
	}
	s.MerkleTree[s.MerkleTreeHeight] = mtList
	for {
		treeLength := len(mtList)
		s.MerkleTreeHeight++
		if treeLength <= 1 {
			break
		} else if treeLength % 2 == 1 {
			mtList = append(mtList, mtList[treeLength - 1])
			treeLength++
		}
		var tmpMtList []string
		for x := 0; x < treeLength/2; x++{
			idx := x * 2
			hashArg := []string{mtList[idx], mtList[idx+1]}
			mkHash := common.ComputeSHA256(hashArg)
			tmpMtList = append(tmpMtList, mkHash)
		}
		mtList = tmpMtList
		s.MerkleTree[s.MerkleTreeHeight] = mtList
	}
	if len(mtList) == 1 {
		s.Header.MerkleTreeRootHash = mtList[0]
	}
}

func (s *Block) MakeMerklePath(){

}



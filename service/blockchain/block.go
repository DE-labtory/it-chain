package blockchain

import (
	"time"
	"errors"
)

type Status int

const (
	Status_BLOCK_UNCONFIRMED Status = 0
	Status_BLOCK_CONFIRMED   Status = 1

	TRANSACTION_SIZE int = 4096
)

type Block struct {
	Header          	*BlockHeader
	MerkleTree      	[2*TRANSACTION_SIZE - 1][32]uint8
	TransactionCount	int
	Transactions	 	[TRANSACTION_SIZE]*Transaction
}

func (s *Block) Reset() { *s = Block{} }

type BlockHeader struct {
	Number             uint64
	PreviousHash       string
	DataHash           string
	Version            string
	MerkleTreeRootHash string
	TimeStamp          time.Time
	BlockHeight        int
	BlockStatus        Status
	CreatedPeerID      string
	Signature          []byte
	PeerId             string
}

func (s *BlockHeader) Reset() { *s = BlockHeader{} }

func (s *Block) PutTranscation(tx Transaction) bool{
	if tx.TransactionStatus == Status_BLOCK_UNCONFIRMED{
		if tx.Validate(){
			tx.TransactionStatus = Status_BLOCK_CONFIRMED
		} else {
			return false
		}
	}
	s.Transactions[s.TransactionCount] = &tx
	s.MerkleTree[TRANSACTION_SIZE - 1 + s.TransactionCount] = tx.TransactionHash
	s.TransactionCount++
	return true
}

func (s Block) FindTransactionIndex(hash [32]uint8) (idx int, err error){
	for idx = 0; idx < s.TransactionCount; idx++{
		if hash == s.Transactions[idx].TransactionHash{
			return idx, nil
		}
	}
	return -1, errors.New("wrong Hash")
}

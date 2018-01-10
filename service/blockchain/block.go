package blockchain

import (
	"time"
	"errors"
)

type Status int32

const (
	Status_UNCONFIRMED Status = 0
	Status_CONFIRMED   Status = 1

	TRANSACTION_SIZE int = 4096
)

type Block struct {
	header          	*BlockHeader
	merkleTree      	[2*TRANSACTION_SIZE - 1][32]uint8
	transactionCount	int
	transactions	 	[TRANSACTION_SIZE]*Transaction
}

func (s *Block) Reset() { *s = Block{} }

type BlockHeader struct {
	number             uint64
	previousHash       string
	dataHash           string
	version            string
	merkleTreeRootHash string
	timeStamp          time.Time
	blockHeight        int
	blockStatus        Status
	createdPeerID      string
	signature          []byte
}

func (s *BlockHeader) Reset() { *s = BlockHeader{} }

func (s *Block) PutTranscation(tx Transaction) bool{
	if tx.transactionStatus == Status_UNCONFIRMED{
		if tx.Validate(){
			tx.transactionStatus = Status_CONFIRMED
		} else {
			return false
		}
	}
	s.transactions[s.transactionCount] = &tx
	s.merkleTree[TRANSACTION_SIZE - 1 + s.transactionCount] = tx.transactionHash
	s.transactionCount++
	return true
}

func (s Block) FindTransactionIndex(hash [32]uint8) (idx int, err error){
	for idx = 0; idx < s.transactionCount; idx++{
		if hash == s.transactions[idx].transactionHash{
			return idx, nil
		}
	}
	return -1, errors.New("wrong Hash")
}

package blockchain

import (
	"testing"
	"fmt"
)

func TestCreateNewTransactionTest(t *testing.T){

}

func TestTxTest(t *testing.T) {
	blk := Block{}
	tx := Transaction{}
	tx.TxData = &TxData{"ss","dd", Params{3, "func", []string{"10", "6"}}, "ff"}
	tx.GenerateHash()

	if blk.PutTranscation(tx) == true{
		idx, err := blk.FindTransactionIndex(tx.TransactionHash)
		fmt.Println(idx, err)
		fmt.Printf("%x\n", blk.Transactions[idx].TransactionHash)
	}

}
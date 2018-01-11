package blockchain

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"fmt"
)

func TestCreateNewTransactionTest(t *testing.T){

}

func TestTxTest(t *testing.T) {
	blk := Block{}
	tx := Transaction{}
	tx.txData = &TxData{"ss","dd", Params{3, "func", []string{"10", "6"}}, "ff"}
	tx.GenerateHash()

	if blk.PutTranscation(tx) == true{
		idx, err := blk.FindTransactionIndex(tx.transactionHash)
		fmt.Println(idx, err)
		fmt.Printf("%x\n", blk.transactions[idx].transactionHash)
	}

}
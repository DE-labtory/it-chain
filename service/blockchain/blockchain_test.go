package blockchain

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"fmt"
)

func TestCreateNewBlockChainTest(t *testing.T){

	var blockChains = CreateNewBlockChain(defaultChannelName,defaultPeerId)

	assert.Equal(t,0,len(blockChains.blocks))
	assert.Equal(t,defaultPeerId,blockChains.header.peerID)
	assert.Equal(t,defaultChannelName,blockChains.header.channelName)


}

func TestTxTest(t *testing.T) {
	blk := Block{}
	tx := Transaction{}
	tx.txData = TxData{"ss","dd", Params{3, "func", []string{"10", "6"}}, "ff"}
	tx.GenerateHash()

	if blk.PutTranscation(tx) == true{
		idx, err := blk.FindTransactionIndex(tx.transactionHash)
		fmt.Println(idx, err)
		fmt.Printf("%x\n", blk.transactions[idx].transactionHash)
	}

}
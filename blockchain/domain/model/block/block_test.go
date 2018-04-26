package block

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	tx "github.com/it-chain/yggdrasill/transaction"

)

func TestCreateGenesisBlock(t *testing.T) {
	GenesisBlock, err := CreateGenesisBlock()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), GenesisBlock.Header.Height)
	assert.Equal(t, "", GenesisBlock.Header.PreviousHash)
	assert.Equal(t, "", GenesisBlock.Header.Version)
	assert.Equal(t, "", GenesisBlock.Header.MerkleTreeRootHash)
	assert.Equal(t, time.Now().String()[:19], GenesisBlock.Header.TimeStamp.String()[:19])
	assert.Equal(t, "Genesis", GenesisBlock.Header.CreatorID)
	assert.Equal(t, make([]byte, 0), GenesisBlock.Header.Signature)
	assert.Equal(t, "", GenesisBlock.Header.BlockHash)
	assert.Equal(t, 0, GenesisBlock.Header.MerkleTreeHeight)
	assert.Equal(t, 0, GenesisBlock.Header.TransactionCount)
	assert.Equal(t, make([][]byte, 0), GenesisBlock.Proof)
	assert.Equal(t, make([]*tx.DefaultTransaction, 0), GenesisBlock.Transactions)
}
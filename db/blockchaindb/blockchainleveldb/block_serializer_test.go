package blockchainleveldb

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/service/blockchain"
)

func Test_Serialize(t *testing.T) {
	block := &blockchain.Block{}

	serialized, err := Serialize(block)
	assert.NoError(t, err)

	deserialized, err := Deserialize(serialized)
	assert.NoError(t, err)
	assert.Equal(t, block, deserialized)
}
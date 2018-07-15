package api_gateway_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/api_gateway"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/stretchr/testify/assert"
)

func TestBlockPoolRepositoryImpl_AddCreatedBlock(t *testing.T) {
	bpr := api_gateway.NewBlockpoolRepositoryImpl()

	// when
	block1 := &blockchain.DefaultBlock{
		Seal:   []byte{0x1},
		Height: blockchain.BlockHeight(1),
	}
	// when
	bpr.AddCreatedBlock(block1)

	// then
	assert.Equal(t, 1, len(bpr.Blocks))
	assert.Equal(t, []byte{0x1}, bpr.Blocks[0].GetSeal())
}

func TestBlockPoolRepositoryImpl_GetStagedBlockByHeight(t *testing.T) {
	bpr := api_gateway.NewBlockpoolRepositoryImpl()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
	})

	// when
	block, err := bpr.GetStagedBlockByHeight(1)

	// then
	assert.Equal(t, err, nil)
	assert.Equal(t, uint64(1), block.GetHeight())
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x1}, block.GetPrevSeal())

	// when
	_, err2 := bpr.GetStagedBlockByHeight(133)

	// then
	assert.Equal(t, err2, api_gateway.ErrNoStagedBlock)
}

func TestBlockPoolRepositoryImpl_GetStagedBlockById(t *testing.T) {
	bpr := api_gateway.NewBlockpoolRepositoryImpl()

	// when
	bpr.Blocks = append(bpr.Blocks, &blockchain.DefaultBlock{
		Seal:     []byte{0x1},
		PrevSeal: []byte{0x1},
		Height:   uint64(1),
	})

	// when
	block, err := bpr.GetStagedBlockById(string([]byte{0x1}))

	// then
	assert.Equal(t, err, nil)
	assert.Equal(t, uint64(1), block.GetHeight())
	assert.Equal(t, []byte{0x1}, block.GetSeal())
	assert.Equal(t, []byte{0x1}, block.GetPrevSeal())

	// when
	_, err2 := bpr.GetStagedBlockById(string([]byte{0x2}))

	// then
	assert.Equal(t, err2, api_gateway.ErrNoStagedBlock)
}

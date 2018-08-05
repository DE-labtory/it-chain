package mem_test

import (
	"os"
	"testing"

	"github.com/it-chain/engine/api_gateway/test/mock"
	"github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/stretchr/testify/assert"
)

func TestBlockRepositoryImpl_FindLast(t *testing.T) {

	dbPath := "./.db"

	// when
	br, err := mem.NewBlockRepositoryImpl(dbPath)

	// then
	assert.Equal(t, nil, err)
	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = br.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
	err = br.AddBlock(block2)
	// then
	assert.NoError(t, err)

	// when
	block3, err := br.FindLast()
	// then
	assert.NoError(t, err)
	assert.Equal(t, block2.GetSeal(), block3.GetSeal())
	assert.Equal(t, uint64(1), block3.GetHeight())

	// when
	AllBlock, err4 := br.FindAll()

	// then
	assert.NoError(t, err4)
	assert.Equal(t, 2, len(AllBlock))

}

package blockchain_test

// ToDo: 삭제 제안 - junk_sound
//func TestBlockPoolModel(t *testing.T) {
//	pool := blockchain.NewBlockPool()
//
//	block1 := &blockchain.DefaultBlock{
//		Height: blockchain.BlockHeight(2),
//	}
//
//	// When
//	pool.Add(block1)
//
//	// Then
//	assert.Equal(t, uint64(2), pool.Get(blockchain.BlockHeight(2)).GetHeight())
//
//	// When
//	block2 := &blockchain.DefaultBlock{
//		Height: blockchain.BlockHeight(2),
//	}
//	pool.Delete(block2)
//
//	// Then
//	assert.Equal(t, nil, pool.Get(blockchain.BlockHeight(2)))
//
//	// when
//	aggregateID := pool.GetID()
//
//	// Then
//	assert.Equal(t, blockchain.BLOCK_POOL_AID, aggregateID)
//}
//
//func TestBlockPoolModel_On(t *testing.T) {
//	pool := blockchain.NewBlockPool()
//
//	event1 := &blockchain.BlockAddToPoolEvent{
//		Height: 1,
//		Seal:   []byte{0x1},
//	}
//	// when
//	err := pool.On(event1)
//	// then
//	assert.Equal(t, nil, err)
//	assert.Equal(t, blockchain.BlockHeight(1), pool.Pool[blockchain.BlockHeight(1)].GetHeight())
//	assert.Equal(t, []byte{0x1}, pool.Pool[blockchain.BlockHeight(1)].GetSeal())
//
//	event2 := &blockchain.BlockAddToPoolEvent{
//		Height: 2,
//		Seal:   []byte{0x2},
//	}
//	// when
//	err2 := pool.On(event2)
//	// then
//	assert.Equal(t, nil, err2)
//	assert.Equal(t, blockchain.BlockHeight(2), pool.Pool[blockchain.BlockHeight(2)].GetHeight())
//	assert.Equal(t, []byte{0x2}, pool.Pool[blockchain.BlockHeight(2)].GetSeal())
//
//	// Same height with event1, but different seal
//	event3 := &blockchain.BlockAddToPoolEvent{
//		Height: 1,
//		Seal:   []byte{0x3},
//	}
//	// when
//	err3 := pool.On(event3)
//	// then
//	assert.Equal(t, nil, err3)
//	assert.Equal(t, blockchain.BlockHeight(1), pool.Pool[blockchain.BlockHeight(1)].GetHeight())
//	assert.Equal(t, []byte{0x3}, pool.Pool[blockchain.BlockHeight(1)].GetSeal())
//
//}

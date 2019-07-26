/*
 * Copyright 2018 DE-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package repo

import (
	"os"
	"testing"

	"github.com/DE-labtory/it-chain/api_gateway/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestBlockRepositoryImpl_NewBlockRepository(t *testing.T) {

	dbPath := "./.db"

	// when
	br, err := NewBlockRepository(dbPath)

	// then
	assert.Equal(t, nil, err)
	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

}

func TestBlockRepositoryImpl_FindBySeal(t *testing.T) {

	dbPath := "./.db"

	// when
	br, err := NewBlockRepository(dbPath)

	// then
	assert.Equal(t, nil, err)
	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block := mock.GetNewBlock([]byte("genesis"), 0)
	err = br.AddBlock(block)
	// then
	assert.NoError(t, err)

	// when
	seal := block.GetSeal()
	blockFindbySeal, err := br.FindBySeal(seal)

	// then
	assert.NoError(t, err)
	assert.Equal(t, block.GetSeal(), blockFindbySeal.GetSeal())

	// when
	block2 := mock.GetNewBlock([]byte("genesis"), 0)
	sealError := block2.GetSeal()
	_, err = br.FindBySeal(sealError)

	// then
	assert.Error(t, err)

}

func TestBlockRepositoryImpl_FindLastAndFindAll(t *testing.T) {

	dbPath := "./.db"

	// when
	br, err := NewBlockRepository(dbPath)

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

func TestBlockRepositoryImpl_FindByHeight(t *testing.T) {

	dbPath := "./.db"

	// when
	br, err := NewBlockRepository(dbPath)

	// then
	assert.Equal(t, nil, err)
	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	blockHeight0 := mock.GetNewBlock([]byte("genesis"), 0)
	err = br.AddBlock(blockHeight0)
	// then
	assert.NoError(t, err)

	// when
	blockFind, err := br.FindByHeight(uint64(0))

	// then
	assert.NoError(t, err)
	assert.Equal(t, blockHeight0.GetSeal(), blockFind.GetSeal())
	assert.Equal(t, blockHeight0.GetHeight(), blockFind.GetHeight())

	// when

	_, err = br.FindByHeight(uint64(1))
	assert.Error(t, err)

}

func TestBlockRepositoryImpl_Save(t *testing.T) {
	dbPath := "./.db"

	// when
	br, err := NewBlockRepository(dbPath)

	// then
	assert.Equal(t, nil, err)
	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block := mock.GetNewBlock([]byte("genesis"), 0)
	err = br.Save(*block)

	//then
	assert.NoError(t, err)
}

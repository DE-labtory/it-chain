/*
 * Copyright 2018 It-chain
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

package kvstore_test

import (
	"reflect"
	"testing"

	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/engine/ivm/infra/kvstore"
	"github.com/stretchr/testify/assert"
)

func TestNewLevelDBStateRepository(t *testing.T) {
	repo := kvstore.NewLevelDBStateRepository()

	assert.Equal(t, reflect.TypeOf(repo), reflect.TypeOf(&kvstore.LevelDBStateRepository{}))
	repo.Close()
}

func TestLevelDBStateRepository_Apply(t *testing.T) {

	//apply success case
	repo := kvstore.NewLevelDBStateRepository()

	writeList := make([]ivm.Write, 0)
	writeList = append(writeList, ivm.Write{
		Key:   []byte("1"),
		Value: []byte("1"),
	})
	writeList = append(writeList, ivm.Write{
		Key:   []byte("2"),
		Value: []byte("2"),
	})

	repo.Apply(writeList)
	v1, _ := repo.Get([]byte("1"))
	v2, _ := repo.Get([]byte("2"))

	assert.Equal(t, v1, []byte("1"))
	assert.Equal(t, v2, []byte("2"))
	repo.Close()

	//	apply fail case
	repo2 := kvstore.NewLevelDBStateRepository()

	writeList2 := make([]ivm.Write, 0)
	writeList2 = append(writeList2, ivm.Write{
		Key:   []byte("1"),
		Value: []byte("1"),
	})
	writeList2 = append(writeList2, ivm.Write{
		Key:   []byte("2"),
		Value: []byte("2"),
	})

	repo2.Apply(writeList)
	repo2.Close()
	v3, _ := repo2.Get([]byte("1"))
	v4, _ := repo2.Get([]byte("2"))

	assert.Equal(t, len(v3), 0)
	assert.Equal(t, len(v4), 0)
}

func TestLevelDBStateRepository_Get(t *testing.T) {
	repo := kvstore.NewLevelDBStateRepository()

	writeList := make([]ivm.Write, 0)
	writeList = append(writeList, ivm.Write{
		Key:   []byte("1"),
		Value: []byte("1"),
	})

	repo.Apply(writeList)
	v, _ := repo.Get([]byte("1"))

	assert.Equal(t, v, []byte("1"))
	repo.Close()

}

func TestLevelDBStateRepository_Close(t *testing.T) {
	repo := kvstore.NewLevelDBStateRepository()
	repo.Close()

	_, err := repo.Get([]byte("1"))

	assert.NotNil(t, err)
}

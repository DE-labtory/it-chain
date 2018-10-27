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

package kvstore

import (
	"sync"

	"os"

	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/iLogger"
	"github.com/it-chain/leveldb-wrapper"
)

type LevelDBStateRepository struct {
	mux sync.RWMutex
	db  *leveldbwrapper.DB
}

func NewLevelDBStateRepository() *LevelDBStateRepository {

	path := os.Getenv("GOPATH") + "/src/github.com/it-chain/engine/state-db"
	os.Remove(path)
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()

	return &LevelDBStateRepository{
		mux: sync.RWMutex{},
		db:  db,
	}
}

func (lr *LevelDBStateRepository) Apply(writeList []ivm.Write) error {
	lr.mux.Lock()
	defer lr.mux.Unlock()

	KVs := make(map[string][]byte, 0)

	for _, write := range writeList {
		if write.Delete {
			KVs[string(write.Key)] = nil

		}
		KVs[string(write.Key)] = write.Value
	}

	iLogger.Infof(nil, "[ivm] Applying to level DB... KVs: %v", KVs)
	err := lr.db.WriteBatch(KVs, true)

	if err != nil {
		iLogger.Infof(nil, "[ivm] Error occurred while Applying to level DB ")
		return err
	}

	return nil
}

func (lr *LevelDBStateRepository) Get(value []byte) ([]byte, error) {
	lr.mux.Lock()
	defer lr.mux.Unlock()

	return lr.db.Get(value)
}

func (lr *LevelDBStateRepository) GetDB() *leveldbwrapper.DB {
	return lr.db
}

func (lr *LevelDBStateRepository) Close() {
	lr.mux.Lock()
	defer lr.mux.Unlock()

	lr.db.Close()
}

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

package api_gateway

import (
	"sync"

	"github.com/it-chain/leveldb-wrapper/key_value_db"
)

type DBHandle struct {
	dbName string
	db     key_value_db.KeyValueDB
}

type DBProvider struct {
	db        key_value_db.KeyValueDB
	mux       sync.Mutex
	dbHandles map[string]*DBHandle
}

func CreateNewDBProvider(keyValueDB key_value_db.KeyValueDB) *DBProvider {
	keyValueDB.Open()
	return &DBProvider{keyValueDB, sync.Mutex{}, make(map[string]*DBHandle)}
}

func (p *DBProvider) Close() {
	p.db.Close()
}

func (p *DBProvider) GetDBHandle(dbName string) *DBHandle {
	p.mux.Lock()
	defer p.mux.Unlock()

	dbHandle := p.dbHandles[dbName]
	if dbHandle == nil {
		dbHandle = &DBHandle{dbName, p.db}
		p.dbHandles[dbName] = dbHandle
	}

	return dbHandle
}

func (h *DBHandle) Get(key []byte) ([]byte, error) {
	return h.db.Get(dbKey(h.dbName, key))
}

func (h *DBHandle) Put(key []byte, value []byte, sync bool) error {
	return h.db.Put(dbKey(h.dbName, key), value, sync)
}

func (h *DBHandle) Delete(key []byte, sync bool) error {
	return h.db.Delete(dbKey(h.dbName, key), sync)
}

func (h *DBHandle) WriteBatch(KVs map[string][]byte, sync bool) error {
	return h.db.WriteBatch(KVs, sync)
}

func (h *DBHandle) GetIteratorWithPrefix() key_value_db.KeyValueDBIterator {
	return h.db.GetIteratorWithPrefix([]byte(h.dbName + "_"))
}

func (h *DBHandle) Snapshot() (map[string][]byte, error) {
	return h.db.Snapshot()
}

func dbKey(dbName string, key []byte) []byte {
	dbName = dbName + "_"
	return append([]byte(dbName), key...)
}

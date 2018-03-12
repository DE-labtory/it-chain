/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package leveldbhelper

import (
	"sync"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/it-chain/it-chain-Engine/common"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type dbState int32

const (
	closed dbState = iota
	opened
)

var logger = common.GetLogger("leveldb_helper.go")

type DB struct {
	levelDBPath     string
	db              *leveldb.DB
	dbState         dbState
	mux             sync.Mutex

	readOpts        *opt.ReadOptions
	writeOptsNoSync *opt.WriteOptions
	writeOptsSync   *opt.WriteOptions
}

//tested
func CreateNewDB(levelDbPath string) *DB {
	readOpts := &opt.ReadOptions{}
	writeOptsNoSync := &opt.WriteOptions{false,false}
	writeOptsSync := &opt.WriteOptions{}
	writeOptsSync.Sync = true

	return &DB{
		levelDBPath:     levelDbPath,
		dbState:         closed,
		readOpts:        readOpts,
		writeOptsNoSync: writeOptsNoSync,
		writeOptsSync:   writeOptsSync,
	}
}

//tested
func (db *DB) Open(){
	db.mux.Lock()
	defer db.mux.Unlock()

	if db.dbState == opened{
		return
	}

	dbOpts := &opt.Options{}
	dbPath := db.levelDBPath
	var err error
	var dirEmpty bool

	if err = common.CreateDirIfMissing(dbPath); err != nil {
		panic(fmt.Sprintf("Error while trying to create dir if missing: %s", err))
	}

	dirEmpty ,err = common.DirEmpty(dbPath)

	dbOpts.ErrorIfMissing = !dirEmpty

	if db.db, err = leveldb.OpenFile(dbPath, dbOpts); err != nil {
		panic(fmt.Sprintf("Error while trying to open DB: %s", err))
	}
	db.dbState = opened
}

func (db *DB) Close(){
	db.mux.Lock()
	defer db.mux.Unlock()

	if db.dbState == closed {
		return
	}
	if err := db.db.Close(); err != nil {
		panic(fmt.Sprintf("Error while trying to close DB: %s", err))
	}
	db.dbState = closed
}

func (db *DB) Get(key []byte) ([]byte, error) {

	value, err := db.db.Get(key, db.readOpts)

	if err == leveldb.ErrNotFound {
		value = nil
		err = nil
	}

	if err != nil {
		logger.Errorf("Error while trying to retrieve key [%#v]: %s", key, err)
		return nil, err
	}
	return value, nil
}


//sync 옵션은 write buffer를 file io와 항상 동기화 시키는지 여부
//sync false면 시스템 crash발생시 정보의 손실 가능성이 있다.
func (db *DB) Put(key []byte, value []byte, sync bool) error {
	wo := db.writeOptsNoSync

	if sync {
		wo = db.writeOptsSync
	}

	err := db.db.Put(key, value, wo)
	if err != nil {
		logger.Errorf("Error while trying to write key [%#v]", key)
		return err
	}
	return nil
}

// Delete deletes the given key
func (db *DB) Delete(key []byte, sync bool) error {
	wo := db.writeOptsNoSync
	if sync {
		wo = db.writeOptsSync
	}
	err := db.db.Delete(key, wo)
	if err != nil {
		logger.Errorf("Error while trying to delete key [%#v]", key)
		return err
	}
	return nil
}

// WriteBatch writes a batch
func (db *DB) WriteBatch(batch *leveldb.Batch, sync bool) error {
	wo := db.writeOptsNoSync
	if sync {
		wo = db.writeOptsSync
	}
	if err := db.db.Write(batch, wo); err != nil {
		return err
	}
	return nil
}

func (db *DB) GetIterator(startKey []byte, endKey []byte) iterator.Iterator {
	return db.db.NewIterator(&util.Range{Start: startKey, Limit: endKey}, db.readOpts)
}

func (db *DB) GetIteratorWithPrefix(prefix []byte) iterator.Iterator {
	return db.db.NewIterator(util.BytesPrefix(prefix), db.readOpts)
}

func (db *DB) Snapshot() (map[string][]byte, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	snap, err := db.db.GetSnapshot()

	if err != nil {
		logger.Error("Error while taking snapshot")
		return nil, err
	}

	data := make(map[string][]byte)
	iter := snap.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		data[string(key)] = append([]byte{}, val...)
	}
	iter.Release()
	if iter.Error() != nil {
		logger.Error("Error with iterator")
		return nil, iter.Error()
	}
	return data, nil
}
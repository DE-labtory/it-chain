package leveldbwrapper

import (
	"sync"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"errors"
	"strings"
	"os"
	"path"
	"io"
	"github.com/it-chain/leveldb-wrapper/key_value_db"
)

type dbState int32

const (
	closed dbState = iota
	opened
)

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

	if err = createDirIfMissing(dbPath); err != nil {
		panic(fmt.Sprintf("Error while trying to create dir if missing: %s", err))
	}

	dirEmpty ,err = isdirEmpty(dbPath)

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
		return nil, errors.New("Error while trying to retrieve key "+string(key)+":"+err.Error())
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
		return errors.New("Error while trying to retrieve key "+string(key)+":"+err.Error())
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
		return errors.New("Error while trying to retrieve key "+string(key)+":"+err.Error())
	}
	return nil
}


// WriteBatch writes a batch
func (db *DB) WriteBatch(KVs map[string][]byte, sync bool) error {
	batch := &leveldb.Batch{}
	for k, v := range KVs {
		key := []byte(k)
		if v == nil {
			batch.Delete(key)
		} else {
			batch.Put(key, v)
		}
	}

	return db.writeBatch(batch, sync)
}

func (db *DB) writeBatch(batch *leveldb.Batch, sync bool) error {
	wo := db.writeOptsNoSync
	if sync {
		wo = db.writeOptsSync
	}
	if err := db.db.Write(batch, wo); err != nil {
		return err
	}
	return nil
}

func (db *DB) GetIterator(startKey []byte, endKey []byte) key_value_db.KeyValueDBIterator {
	return db.db.NewIterator(&util.Range{Start: startKey, Limit: endKey}, db.readOpts)
}

func (db *DB) GetIteratorWithPrefix(prefix []byte) key_value_db.KeyValueDBIterator {
	return db.db.NewIterator(util.BytesPrefix(prefix), db.readOpts)
}

func (db *DB) Snapshot() (map[string][]byte, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	snap, err := db.db.GetSnapshot()

	if err != nil {
		return nil, errors.New("Error while taking snapshot:"+err.Error())
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
		return nil, iter.Error()
	}
	return data, nil
}

func createDirIfMissing(dirPath string) (error){

	if !strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath + "/"
	}

	err := os.MkdirAll(path.Dir(dirPath), 0755)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while creating dir [%s]", dirPath))
	}

	return nil
}

// DirEmpty returns true if the dir at dirPath is empty
func isdirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Error while opening dir [%s]: %s", dirPath, err))
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
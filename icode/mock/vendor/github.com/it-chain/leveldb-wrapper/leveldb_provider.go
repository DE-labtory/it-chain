package leveldbwrapper

import (
	"sync"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/it-chain/leveldb-wrapper/key_value_db"
)

type DBHandle struct {
	dbName	string
	db		*DB
}

type DBProvider struct {
	db			*DB
	mux			sync.Mutex
	dbHandles	map[string]*DBHandle
}

func CreateNewDBProvider(levelDbPath string) *DBProvider {
	db := CreateNewDB(levelDbPath)
	db.Open()
	return &DBProvider{db, sync.Mutex{}, make(map[string]*DBHandle)}
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
	batch := &leveldb.Batch{}
	for k, v := range KVs {
		key := dbKey(h.dbName, []byte(k))
		if v == nil {
			batch.Delete(key)
		} else {
			batch.Put(key, v)
		}
	}
	return h.db.writeBatch(batch, sync)
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
package leveldbhelper

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestCreateNewDBProvider(t *testing.T) {
	path := "./test_db_path"
	dbProvider := CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
		os.RemoveAll(path)
	}()

	assert.Equal(t, dbProvider.db.dbState, opened)
}

func TestDBProvider_Close(t *testing.T) {
	path := "./test_db_path"
	dbProvider := CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
		os.RemoveAll(path)
	}()

	dbProvider.Close()

	assert.Equal(t, dbProvider.db.dbState, closed)
}

func TestDBProvider_GetDBHandle(t *testing.T) {
	path := "./test_db_path"
	dbProvider := CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
		os.RemoveAll(path)
	}()

	dbName1 := "db1"
	dbName2 := "db2"
	dbName3 := "db3"

	dbHandle1 := dbProvider.GetDBHandle(dbName1)
	dbHandle2 := dbProvider.GetDBHandle(dbName2)
	dbHandle3 := dbProvider.GetDBHandle(dbName3)

	assert.Equal(t, dbName1, dbHandle1.dbName)
	assert.Equal(t, dbName2, dbHandle2.dbName)
	assert.Equal(t, dbName3, dbHandle3.dbName)

	assert.Equal(t, dbProvider.db, dbHandle1.db)
	assert.Equal(t, dbProvider.db, dbHandle2.db)
	assert.Equal(t, dbProvider.db, dbHandle3.db)
}

func TestDBHandle_Get(t *testing.T) {
	path := "./test_db_path"
	dbProvider := CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
		os.RemoveAll(path)
	}()

	dbName1 := "db1"
	dbName2 := "db2"
	dbName3 := "db3"

	dbHandle1 := dbProvider.GetDBHandle(dbName1)
	dbHandle2 := dbProvider.GetDBHandle(dbName2)
	dbHandle3 := dbProvider.GetDBHandle(dbName3)

	dbHandle1.db.db.Put(dbKey(dbName1, []byte("key")), []byte("value1"), nil)
	dbHandle2.db.db.Put(dbKey(dbName2, []byte("key")), []byte("value2"), nil)
	dbHandle3.db.db.Put(dbKey(dbName3, []byte("key")), []byte("value3"), nil)

	val1, err := dbHandle1.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, string(val1), "value1")

	val2, err := dbHandle2.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, string(val2), "value2")

	val3, err := dbHandle3.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, string(val3), "value3")
}

func TestDBHandle_Put(t *testing.T) {
	path := "./test_db_path"
	dbProvider := CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
		os.RemoveAll(path)
	}()

	dbName1 := "db1"
	dbName2 := "db2"
	dbName3 := "db3"

	dbHandle1 := dbProvider.GetDBHandle(dbName1)
	dbHandle2 := dbProvider.GetDBHandle(dbName2)
	dbHandle3 := dbProvider.GetDBHandle(dbName3)

	dbHandle1.Put([]byte("key"), []byte("value1"), true)
	dbHandle2.Put([]byte("key"), []byte("value2"), true)
	dbHandle3.Put([]byte("key"), []byte("value3"), true)

	val1, err := dbHandle1.db.db.Get(dbKey(dbName1, []byte("key")), nil)
	assert.NoError(t, err)
	assert.Equal(t, string(val1), "value1")

	val2, err := dbHandle2.db.db.Get(dbKey(dbName2, []byte("key")), nil)
	assert.NoError(t, err)
	assert.Equal(t, string(val2), "value2")

	val3, err := dbHandle3.db.db.Get(dbKey(dbName3, []byte("key")), nil)
	assert.NoError(t, err)
	assert.Equal(t, string(val3), "value3")
}

func TestDBHandle_Delete(t *testing.T) {
	path := "./test_db_path"
	dbProvider := CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
		os.RemoveAll(path)
	}()

	dbName := "db"
	dbHandle := dbProvider.GetDBHandle(dbName)

	dbHandle.db.db.Put(dbKey(dbName, []byte("test")), []byte("test"), nil)
	err := dbHandle.Delete([]byte("test"), true)
	assert.NoError(t, err)

	_, err = dbHandle.db.db.Get(dbKey(dbName, []byte("test")), nil)
	assert.Equal(t, err, leveldb.ErrNotFound)
}
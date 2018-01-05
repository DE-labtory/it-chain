package blockchainleveldb

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"github.com/syndtr/goleveldb/leveldb"
)

func getDB(path string) *DB{

	readOpts := &opt.ReadOptions{}
	writeOptsNoSync := &opt.WriteOptions{}
	writeOptsSync := &opt.WriteOptions{}
	writeOptsSync.Sync = true

	return &DB{
		levelDBPath:     path,
		dbState:         closed,
		readOpts:        readOpts,
		writeOptsNoSync: writeOptsNoSync,
		writeOptsSync:   writeOptsSync,
	}
}

func TestCreateNewDB(t *testing.T){

	//when
	path := "./test_db_path"

	//then
	db := CreateNewDB(path)

	assert.Equal(t,db.dbState,closed)
	assert.Equal(t,db.levelDBPath,path)
}

func TestDB_Open(t *testing.T) {

	//when
	path := "./test_db_path/"
	db := getDB(path)

	//then
	db.Open()

	//result
	assert.Equal(t,db.dbState,opened)
	assert.NotNil(t,db.db)

	//clean up
	os.RemoveAll(path)
}

func TestDB_Get_ExistValue(t *testing.T) {
	//when
	path := "./test_db_path/"
	db := getDB(path)
	db.Open()
	err := db.db.Put([]byte("jun"),[]byte("jun"),nil)
	if err != nil{
		t.Error("get error!")
	}

	//then
	byteValue, err := db.Get([]byte("jun"))


	//result
	assert.Equal(t,string(byteValue),"jun")

	//clean up
	os.RemoveAll(path)
}

func TestDB_Get_UNExistValue(t *testing.T) {
	//when
	path := "./test_db_path/"
	db := getDB(path)
	db.Open()
	//err := db.db.Put([]byte("jun"),[]byte("jun"),nil)
	//if err != nil{
	//	t.Error("get error!")
	//}

	//then
	byteValue, err := db.Get([]byte("jun"))

	//result
	assert.Equal(t,[]byte(nil),byteValue)
	assert.Equal(t,nil,err)

	//clean up
	os.RemoveAll(path)
}

//없으면 write 있으면 update로 작용
func TestDB_Put(t *testing.T) {

	//when
	path := "./test_db_path/"
	db := getDB(path)
	db.Open()

	//then
	db.Put([]byte("jun"),[]byte("jun"),true)


	//result
	byteValue, err := db.db.Get([]byte("jun"),nil)
	if err != nil{
		t.Error("get error!")
	}
	assert.Equal(t,string(byteValue),"jun")

	//clean up
	os.RemoveAll(path)
}

//UNExistValue삭제는 아무것도 하지 않음 따라서 existValue만 체크
func TestDB_Delete_ExistValue(t *testing.T) {

	//when
	path := "./test_db_path/"
	db := getDB(path)
	db.Open()
	var err error
	err = db.db.Put([]byte("jun"),[]byte("jun"),nil)
	if err != nil{
		t.Error("put error!")
	}

	//then
	err = db.Delete([]byte("jun"),true)
	if err != nil{
		t.Error("delete error!")
	}

	_, err = db.db.Get([]byte("jun"),nil)

	assert.Equal(t,err,leveldb.ErrNotFound)

	//clean up
	os.RemoveAll(path)

}
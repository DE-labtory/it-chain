package main

import (
	"os"
	"it-chain/db/leveldbhelper"
)

func main() {

	path := "./test"
	dbProvider := leveldbhelper.CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
	}()

	wsDB := "worldStateDB"
	wsDBHandle := dbProvider.GetDBHandle(wsDB)
	wsDBHandle.Put([]byte("key"), []byte("value1"), true)
}
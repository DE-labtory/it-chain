package main

import (
	"it-chain/db/leveldbhelper"
)

func main() {

	path := "./sample/db/test"
	dbProvider := leveldbhelper.CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
	}()

	wsDB := "worldStateDB"
	wsDBHandle := dbProvider.GetDBHandle(wsDB)
	wsDBHandle.Put([]byte("a"), []byte("value1"), true)
}
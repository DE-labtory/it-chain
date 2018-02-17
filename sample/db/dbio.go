package main

import (
	"it-chain/db/leveldbhelper"
)

func main() {

	//path := "./sample/db/test"
	path := "./smartcontract/worldstatedb"
	dbProvider := leveldbhelper.CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
	}()

	wsDB := "worldStateDB"
	wsDBHandle := dbProvider.GetDBHandle(wsDB)
	wsDBHandle.Put([]byte("A"), []byte("12"), true)
}
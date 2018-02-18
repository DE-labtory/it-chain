package main

import (
	"it-chain/db/leveldbhelper"
	"fmt"
)

func main() {

	path := "./sample/db/test"
	//path := "./smartcontract/worldstatedb"
	dbProvider := leveldbhelper.CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
	}()

	wsDB := "worldStateDB"
	wsDBHandle := dbProvider.GetDBHandle(wsDB)
	a, _ := wsDBHandle.Get([]byte("A"))
	wsDBHandle.Put([]byte("A"), []byte("13"), true)

	fmt.Println(string(a))
}
package main

import (
	"it-chain/db/leveldbhelper"
	"fmt"
	"os"
)

func main() {
	GOPATH := os.Getenv("GOPATH")

	//path := "./sample/db/test"
	path := GOPATH + "/src/it-chain/smartcontract/worldstatedb/test"
	dbProvider := leveldbhelper.CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
	}()

	wsDB := "worldStateDB"
	wsDBHandle := dbProvider.GetDBHandle(wsDB)
	a, _ := wsDBHandle.Get([]byte("A"))
	wsDBHandle.Put([]byte("A"), []byte("asdasdasd"), true)

	fmt.Println(string(a))
}
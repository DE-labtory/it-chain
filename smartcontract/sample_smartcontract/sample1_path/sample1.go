package sample1_path

import (
	"it-chain/service/blockchain"
	"it-chain/db/leveldbhelper"
	"os"
	"encoding/json"
	"fmt"
)

type SampleSmartContract struct {

}

func (sc *SampleSmartContract) Init(args []string) error {
	tx := blockchain.Transaction{}
	err := json.Unmarshal([]byte(args[0]), &tx)
	if err != nil {
		return err
	}


	/* Init WorldStateDB
	---------------------*/
	path := "./"
	dbProvider := leveldbhelper.CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
	}()

	wsDB := "worldStateDB"
	wsDBHandle := dbProvider.GetDBHandle(wsDB)
	//wsDBHandle.Put([]byte("key"), []byte("value1"), true)

	if tx.TxData.Method == "Query" {
		sc.Query(tx, wsDBHandle)
	} else if tx.TxData.Method == "Invoke" {
		sc.Invoke(tx, wsDBHandle)
	}

	return nil
}

func (sc *SampleSmartContract) Query(transaction blockchain.Transaction, wsDBHandle *leveldbhelper.DBHandle) {
}

func (sc *SampleSmartContract) Invoke(transaction blockchain.Transaction, wsDBHandle *leveldbhelper.DBHandle) {
	//wsDBHandle.Put([]byte("test"), []byte("success"), true)
	fmt.Println("fun Invoke")
}

func main()  {
	args := os.Args[1:]
	ssc := new(SampleSmartContract)

	ssc.Init(args)
}
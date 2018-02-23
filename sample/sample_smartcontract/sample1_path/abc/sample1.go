package main

import (
	"it-chain/domain"
	//"it-chain/common"
	"it-chain/db/leveldbhelper"
	"encoding/json"
	"os"
	"fmt"
)
//var logger = common.GetLogger("smartcontract")

/*** Set SmartContractResponse ***/
var smartContractResponse = domain.SmartContractResponse{Data: map[string]string{}}

type SampleSmartContract struct {
}

func (sc *SampleSmartContract) Init(tx_string string, dbPath string, args []string) error {
	/*** Set Transaction Struct ***/
	tx := domain.Transaction{}
	err := json.Unmarshal([]byte(tx_string), &tx)
	if err != nil {
		return err
	}

	/*** Init WorldStateDB ***/
	path := "/go/src" + dbPath	// "/go/src/it-chain/smartcontract/worldstatedb/test"
	dbProvider := leveldbhelper.CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
	}()

	wsDB := "worldStateDB"
	wsDBHandle := dbProvider.GetDBHandle(wsDB)

	/*** Mock Data ***/
	//wsDBHandle.Put([]byte("A"), []byte("AAAAAAAA"), true)

	if tx.TxData.Method == domain.Query {
		sc.Query(tx, wsDBHandle)
	} else if tx.TxData.Method == domain.Invoke {
		sc.Invoke(tx, wsDBHandle)
	}

	return nil
}

func (sc *SampleSmartContract) Query(tx domain.Transaction, wsDBHandle *leveldbhelper.DBHandle) {
	if tx.TxData.Params.Function == "getA" {
		getA(wsDBHandle)
	}
}

func (sc *SampleSmartContract) Invoke(tx domain.Transaction, wsDBHandle *leveldbhelper.DBHandle) {
	if tx.TxData.Params.Function == "putA" {
		putA(wsDBHandle)
	}
	return
}

func main()  {
	defer func() {
		if err := recover(); err != nil {
			smartContractResponse.Result = domain.FAIL
			smartContractResponse.Error = "An error occured while running smartcontract!"
			return
		} else {
			smartContractResponse.Result = domain.SUCCESS
		}

		/*** Marshal SmartContractResponse ***/
		out, err := json.Marshal(smartContractResponse)
		if err != nil {
			smartContractResponse.Result = domain.FAIL
			smartContractResponse.Error = "An error occured while marshaling the response!"
			return
		}
		fmt.Println(string(out))
	}()

	tx_string := os.Args[1]
	dbPath := os.Args[2]
	var args []string
	if len(os.Args) > 3 {
		args = os.Args[3:]
	}
	ssc := new(SampleSmartContract)

	ssc.Init(tx_string, dbPath, args)
}


/* Custom Function
 -----------------------*/

func getA(wsDBHandle *leveldbhelper.DBHandle) {
	a, _ := wsDBHandle.Get([]byte("A"))
	smartContractResponse.Data["A"] = string(a)
}

func putA(wsDBHandle *leveldbhelper.DBHandle) {
	wsDBHandle.Put([]byte("A"), []byte("123123"), true)
}
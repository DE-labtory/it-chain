package main

import (
	"it-chain/domain"
	//"it-chain/common"
	"it-chain/db/leveldbhelper"
	"it-chain/smartcontract"
	"encoding/json"
	"os"
	"fmt"
)
//var logger = common.GetLogger("smartcontract")

/*** Set SmartContractResponse ***/
var smartContractResponse = smartcontract.SmartContractResponse{Data: map[string]string{}}

type SampleSmartContract struct {
}

func (sc *SampleSmartContract) Init(args []string) error {
	/*** Set Transaction Struct ***/
	tx := domain.Transaction{}
	err := json.Unmarshal([]byte(args[0]), &tx)
	if err != nil {
		return err
	}

	/*** Init WorldStateDB ***/
	path := "/go/src/worldstatedb"
	dbProvider := leveldbhelper.CreateNewDBProvider(path)
	defer func(){
		dbProvider.Close()
	}()

	wsDB := "worldStateDB"
	wsDBHandle := dbProvider.GetDBHandle(wsDB)

	/*** Mock Data ***/
	wsDBHandle.Put([]byte("A"), []byte("AAAAAAAA"), true)

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
	return
}

func main()  {
	defer func() {
		if err := recover(); err != nil {
			smartContractResponse.Result = smartcontract.FAIL
			smartContractResponse.Error = "An error occured while running smartcontract!"
			return
		} else {
			smartContractResponse.Result = smartcontract.SUCCESS
		}

		/*** Marshal SmartContractResponse ***/
		out, err := json.Marshal(smartContractResponse)
		if err != nil {
			smartContractResponse.Result = smartcontract.FAIL
			smartContractResponse.Error = "An error occured while marshaling the response!"
			return
		}
		fmt.Println(string(out))
	}()

	args := os.Args[1:]
	ssc := new(SampleSmartContract)

	ssc.Init(args)
}


/* Custom Function
 -----------------------*/

func getA(wsDBHandle *leveldbhelper.DBHandle) {
	a, _ := wsDBHandle.Get([]byte("A"))
	smartContractResponse.Data["A"] = string(a)
}
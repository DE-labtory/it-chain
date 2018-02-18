package main

import (
	"it-chain/domain"
	"it-chain/common"
	"it-chain/db/leveldbhelper"
	"it-chain/smartcontract"
	"encoding/json"
	"os"
	"fmt"
	"errors"
)
var logger = common.GetLogger("smartcontract")

/*** Set SmartContractResponse ***/
var smartcontract_response = smartcontract.SmartContractResponse{}

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
	//wsDBHandle.Put([]byte("test"), []byte("value1"), true)

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
			smartcontract_response.Result = smartcontract.FAIL
			smartcontract_response.Error = errors.New("An error occured while running smartcontract!")
		} else {
			smartcontract_response.Result = smartcontract.SUCCESS
		}

		/*** Marshal SmartContractResponse ***/
		out, err := json.Marshal(smartcontract_response)
		if err != nil {
			smartcontract_response.Result = smartcontract.FAIL
			smartcontract_response.Error = errors.New("An error occured while marshaling the response!")
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
	smartcontract_response.Data["A"] = string(a)
}
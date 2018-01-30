package smartcontract

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/service/blockchain"
	"fmt"
	"strconv"
	"time"
)

func TestDeploy_Deploy(t *testing.T) {
	scs := SmartContractService{"hackurity01", map[string]SmartContract{}}
	ContractPath := "junbeomlee/bloom"

	deploy_result, err := scs.Deploy(ContractPath)

	fmt.Println(deploy_result)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	fmt.Println(deploy_result)
	//assert.Equal(t,nil,deploy_result)
}

func TestSmartContractService_Query(t *testing.T) {
	tx := blockchain.CreateNewTransaction(
		strconv.Itoa(1),
		strconv.Itoa(1),
		0,
		time.Now(),
		blockchain.SetTxData(
			"",
			"query",
			blockchain.SetTxMethodParameters(0, "", []string{""}),
			"abc",
		),
	)
	scs := SmartContractService{
		"hackurity01",
		map[string]SmartContract{
			"abc": SmartContract{
				ReposName:         "bloom",
				OriginReposPath:   "junbeomlee/bloom",
				SmartContractPath: "/Users/hackurity/Documents/it-chain/test/abc",
			},
		},
	}

	scs.Query(*tx)

}

func TestSmartContractService_Invoke(t *testing.T) {

}

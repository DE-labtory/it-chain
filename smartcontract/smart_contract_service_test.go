package smartcontract

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/service/blockchain"
	"fmt"
	"strconv"
	"time"
	"path/filepath"
)

func TestDeploy_Deploy(t *testing.T) {
	currentDir, err := filepath.Abs("./")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	scs := SmartContractService{"steve-buzzni", currentDir + "/sample_smartcontract",map[string]SmartContract{}}
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
	currentDir, err := filepath.Abs("./")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	tx := blockchain.CreateNewTransaction(
		strconv.Itoa(1),
		strconv.Itoa(1),
		0,
		time.Now(),
		blockchain.SetTxData(
			"",
			"Invoke",
			blockchain.SetTxMethodParameters(0, "", []string{""}),
			"sample1_path",
		),
	)
	scs := SmartContractService{
		"hackurity01",
		currentDir + "/sample_smartcontract",
		map[string]SmartContract{
			"sample1_path": SmartContract{
				ReposName:         "sample1",
				OriginReposPath:   "sample/sample1",
				SmartContractPath: currentDir + "/sample_smartcontract/sample1_path",
			},
		},
	}

	scs.Query(*tx)

}

func TestSmartContractService_Invoke(t *testing.T) {

}

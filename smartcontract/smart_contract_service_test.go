package smartcontract

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"os"
	"fmt"
	"it-chain/domain"
	"strconv"
	"time"
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
	tx := domain.CreateNewTransaction(
		strconv.Itoa(1),
		strconv.Itoa(1),
		0,
		time.Now(),
		domain.SetTxData(
			"",
			"query",
			domain.SetTxMethodParameters(0, "", []string{""}),
			"abc",
		),
	)
	scs := SmartContractService{
		"hackurity01",
		currentDir + "/sample_smartcontract",
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

func TestSmartContractService_pullAllSmartContracts(t *testing.T) {

	currentDir, err := filepath.Abs("./")
	if err != nil {
		assert.Fail(t, err.Error())
	}

	scs := SmartContractService{
		"yojkim",
	currentDir + "/pull_test_repositories",
	map[string]SmartContract{}}

	scs.pullAllSmartContracts("emperorhan", func(e error) {
		assert.Fail(t, e.Error())
	}, nil)

	scs.pullAllSmartContracts("", func(e error) {
		assert.Error(t, e)
	}, nil)

	defer os.RemoveAll(scs.SmartContractDirPath)

}

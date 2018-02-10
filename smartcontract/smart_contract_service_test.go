package smartcontract

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/domain"
	"fmt"
	"strconv"
	"time"
	"path/filepath"
	"os"
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
			"Invoke",
			domain.SetTxMethodParameters(0, "", []string{""}),
			"abc",
		),
	)
	fmt.Println("tx created")
	scs := SmartContractService{
		"steve-buzzni",
		currentDir + "/sample_smartcontract",
		map[string]SmartContract{
			"abc": SmartContract{
				ReposName:         "sample1",
				OriginReposPath:   "sample1/path",
				SmartContractPath: currentDir + "/sample_smartcontract/sample1_path",
			},
		},
	}

	fmt.Println("scs created")
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

	err = scs.pullAllSmartContracts()
	assert.NoError(t, err)

	defer os.RemoveAll(scs.SmartContractDirPath)

}

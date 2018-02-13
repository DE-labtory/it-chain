package smartcontract

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"it-chain/domain"
	"strconv"
	"time"
	"path/filepath"
	"os"
	"os/exec"
	"bytes"
	"io"
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
				Name:         "sample1",
				OriginReposPath:   "sample1/path",
				SmartContractPath: currentDir + "/sample_smartcontract/sample1_path",
				//SmartContractPath: "/Users/hackurity/go/src/it-chain-smartcontract/sample1_path",
			},
		},
	}

	fmt.Println("scs created")
	scs.Query(*tx)

	defer func() {
		//docker rm $(docker ps -a -f "ancestor=golang:1.9.2-alpine3.6" -q)
		//docker ps -a -f "ancestor=golang:1.9.2-alpine3.6" -q | xargs -I {} docker rm {}
		c1 := exec.Command("docker", "ps", "-a", "-f", "ancestor=golang:1.9.2-alpine3.6", "-q")
		c2 := exec.Command("xargs", "-I", "{}", "docker", "rm", "{}")

		r, w := io.Pipe()
		c1.Stdout = w
		c2.Stdin = r

		var b2 bytes.Buffer
		c2.Stdout = &b2

		c1.Start()
		c2.Start()
		c1.Wait()
		w.Close()
		c2.Wait()
	}()
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

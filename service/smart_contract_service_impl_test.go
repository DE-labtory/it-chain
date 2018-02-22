package service

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

//todo deploy
func TestDeploy_Deploy(t *testing.T) {
	currentDir, err := filepath.Abs("./")
	if err != nil {
		assert.Fail(t, err.Error())
	}

	scs := SmartContractServiceImpl{"steve-buzzni", currentDir + "/sample_smartcontract",map[string]SmartContract{}}
	ContractPath := "junbeomlee/bloom"

	deploy_result, err := scs.Deploy(ContractPath)

	fmt.Println(deploy_result)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	fmt.Println(deploy_result)
	//assert.Equal(t,nil,deploy_result)
}

func TestSmartContractServiceImpl_Invoke(t *testing.T) {

}

func TestSmartContractServiceImpl_ValidateTransactionsOfBlock(t *testing.T) {
	currentDir, err := filepath.Abs("./")
	fmt.Println(currentDir)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	/*** Create Block ***/
	block := domain.CreateNewBlock(nil, "")

	/*** Create Transaction1 ***/
	tx := domain.CreateNewTransaction(
		strconv.Itoa(1),
		strconv.Itoa(1),
		0,
		time.Now(),
		domain.SetTxData(
			"",
			domain.Query,
			domain.SetTxMethodParameters(0, "getA", []string{""}),
			"abc",
		),
	)
	tx.GenerateHash()
	err = block.PutTranscation(tx)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	/*** Create Transaction2 ***/
	tx2 := domain.CreateNewTransaction(
		strconv.Itoa(1),
		strconv.Itoa(1),
		0,
		time.Now(),
		domain.SetTxData(
			"",
			domain.Query,
			domain.SetTxMethodParameters(0, "getA", []string{""}),
			"abc",
		),
	)
	tx2.GenerateHash()
	err = block.PutTranscation(tx2)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	fmt.Println("block created with tx")
	scs := SmartContractServiceImpl{
		"steve-buzzni",
		currentDir + "/sample_smartcontract",
		map[string]SmartContract{
			"abc": SmartContract{
				Name:         "sample1",
				OriginReposPath:   "sample1/path",
				SmartContractPath: currentDir + "/../sample/sample_smartcontract/sample1_path",
				//SmartContractPath: "/Users/hackurity/go/src/it-chain-smartcontract/sample1_path",
			},
		},
	}

	fmt.Println("scs created")
	scs.ValidateTransactionsOfBlock(block)

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

func TestSmartContractServiceImpl_ValidateTransaction(t *testing.T) {
	currentDir, err := filepath.Abs("./")
	fmt.Println(currentDir)
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
			domain.Query,
			domain.SetTxMethodParameters(0, "getA", []string{""}),
			"abc",
		),
	)
	fmt.Println("tx created")
	scs := SmartContractServiceImpl{
		"steve-buzzni",
		currentDir + "/sample_smartcontract",
		map[string]SmartContract{
			"abc": SmartContract{
				Name:         "sample1",
				OriginReposPath:   "sample1/path",
				SmartContractPath: currentDir + "/../sample/sample_smartcontract/sample1_path",
				//SmartContractPath: "/Users/hackurity/go/src/it-chain-smartcontract/sample1_path",
			},
		},
	}

	fmt.Println("scs created")
	scs.ValidateTransaction(tx)

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

func TestSmartContractServiceImpl_pullAllSmartContracts(t *testing.T) {

	currentDir, err := filepath.Abs("./")
	if err != nil {
		assert.Fail(t, err.Error())
	}

	scs := SmartContractServiceImpl{
		"yojkim",
	currentDir + "/pull_test_repositories",
	map[string]SmartContract{}}

	scs.PullAllSmartContracts("emperorhan", func(e error) {
		assert.Fail(t, e.Error())
	}, nil)

	scs.PullAllSmartContracts("", func(e error) {
		assert.Error(t, e)
	}, nil)

	defer os.RemoveAll(scs.SmartContractDirPath)

}

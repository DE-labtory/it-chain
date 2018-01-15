package chaincode

import (
	"errors"
)

const (
	GITHUB_TOKEN string = "75ecdfcf125761352a406bd1d06d208afffb9482"
)

type Chaincode struct {
	SmartContractMap map[string]string
}

func Init() {
}

func (chaincode *Chaincode) Deploy(ContractPath string) (error) {
	_, ok := chaincode.SmartContractMap[ContractPath];
	if ok {
		return errors.New("Already exist smart contract ID")
	}

	_, err := GetRepos(ContractPath)
	if err != nil {
		return errors.New("An error occured while getting repos!")
	}

	githubResponse, err := ForkRepos(ContractPath, GITHUB_TOKEN)
	if err != nil {
		return errors.New("An error occured while forking repos!")
	}

	chaincode.SmartContractMap[ContractPath] = githubResponse.Full_name

	return nil
}

func (chaincode *Chaincode) Query() {

}


func (chaincode *Chaincode) Invoke() {

}
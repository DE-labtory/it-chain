package smartcontract

import (
	"errors"
	"it-chain/service/blockchain"
	"os"
)

const (
	GITHUB_TOKEN string = "75ecdfcf125761352a406bd1d06d208afffb9482"
)

type SmartContract struct {
	RepoName string
	ContractPath string
}

type SmartContractManager struct {
	SmartContractMap map[string]SmartContract
}

func Init() {
}

func (scm *SmartContractManager) Deploy(ContractID string) (error) {
	_, ok := scm.SmartContractMap[ContractID];
	if ok {
		return errors.New("Already exist smart contract ID")
	}

	smartcontract, err := GetRepos(ContractPath)
	if err != nil {
		return errors.New("An error occured while getting repos!")
	}

	githubResponse, err := ForkRepos(ContractPath, GITHUB_TOKEN)
	if err != nil {
		return errors.New("An error occured while forking repos!")
	}

	scm.SmartContractMap[ContractPath] = githubResponse.Full_name

	return nil
}
/***************************************************
 *	1. smartcontract 검사
 *	2. smartcontract -> sc.tar
 *	3. go 버전에 맞는 docker image를 Create
 *	4. sc.tar를 docker container로 복사
 *	5. docker container Start
 *	6. docker에서 smartcontract 실행
 ****************************************************/
func (scm *SmartContractManager) Query(transaction blockchain.Transaction) (error) {

	_, ok := scm.SmartContractMap[transaction.TxData.ContractID];
	if !ok {
		return errors.New("Not exist contract ID")
	}

	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}

	return nil
}


func (smartcontract *SmartContractManager) Invoke() {

}
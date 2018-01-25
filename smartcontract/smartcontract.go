package smartcontract

import (
	"errors"
	"it-chain/service/blockchain"
)

const (
	GITHUB_TOKEN string = "75ecdfcf125761352a406bd1d06d208afffb9482"
)

type SmartContract struct {
	SmartContractMap map[string]string
}

func Init() {
}

func (smartcontract *SmartContract) Deploy(ContractPath string) (error) {
	_, ok := smartcontract.SmartContractMap[ContractPath];
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

	smartcontract.SmartContractMap[ContractPath] = githubResponse.Full_name

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
func (smartcontract *SmartContract) Query(transaction blockchain.Transaction) (error) {

	_, ok := smartcontract.SmartContractMap[transaction.TxData.ContractID];
	if !ok {
		return errors.New("Not exist contract ID")
	}


	return nil
}


func (smartcontract *SmartContract) Invoke() {

}
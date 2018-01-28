package smartcontract

import (
	"errors"
	"it-chain/service/blockchain"
	"strings"
	"fmt"
	"os"
	"time"
)

const (
	GITHUB_TOKEN string = "619ee819fca81e424589a5f10416c2e2e99ed744"

)

type SmartContract struct {
	RepoName string
	ContractPath string
}

type SmartContractService struct {
	GithubID string
	SmartContractMap map[string]SmartContract
}

func Init() {
}

func (scs *SmartContractService) Deploy(ContractPath string) (string, error) {
	origin_repos_name := strings.Split(ContractPath, "/")[1]
	new_repos_name := strings.Replace(ContractPath, "/", "_", -1)

	_, ok := scs.keyByValue(ContractPath)
	if ok {
		// 버전 업데이트
		return "", errors.New("Already exist smart contract ID")
	}

	repos, err := GetRepos(ContractPath)
	if err != nil {
		return "", errors.New("An error occured while getting repos!")
	}
	if repos.Message == "Bad credentials" {
		return "", errors.New("Not Exist Repos!")
	}

	err = CloneRepos(ContractPath, "/Users/hackurity/Documents/it-chain/test")
	if err != nil {
		return "", errors.New("An error occured while cloning repos!")
	}

	fmt.Println(new_repos_name)
	_, err = CreateRepos(new_repos_name, GITHUB_TOKEN)
	if err != nil {
		return "", errors.New(err.Error())//"An error occured while creating repos!")
	}

	err = ChangeRemote(scs.GithubID + "/" + new_repos_name, "/Users/hackurity/Documents/it-chain/test/" + origin_repos_name)
	if err != nil {
		return "", errors.New("An error occured while cloning repos!")
	}

	// 버전 관리를 위한 파일 추가
	now := time.Now().Format("2006-01-02 15:04:05");
	file, err := os.OpenFile("/Users/hackurity/Documents/it-chain/test/" + origin_repos_name + "/version", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return "", errors.New("An error occured while creating or opening file!")
	}

	_, err = file.WriteString("Deployed at " + now + "\n")
	if err != nil {
		return "", errors.New("An error occured while writing file!")
	}
	err = file.Close()
	if err != nil {
		return "", errors.New("An error occured while closing file!")
	}

	err = CommitAndPush("/Users/hackurity/Documents/it-chain/test/" + origin_repos_name, "It-Chain Smart Contract \"" + new_repos_name + "\" Deploy")
	if err != nil {
		return "", errors.New(err.Error())
	}

	githubResponseCommits, err := GetReposCommits(scs.GithubID + "/" + new_repos_name)

	scs.SmartContractMap[githubResponseCommits[0].Sha] = SmartContract{new_repos_name, ContractPath}

	return githubResponseCommits[0].Sha, nil
}
/***************************************************
 *	1. smartcontract 검사
 *	2. smartcontract -> sc.tar
 *	3. go 버전에 맞는 docker image를 Create
 *	4. sc.tar를 docker container로 복사
 *	5. docker container Start
 *	6. docker에서 smartcontract 실행
 ****************************************************/
func (scs *SmartContractService) Query(transaction blockchain.Transaction) (error) {

	_, ok := scs.SmartContractMap[transaction.TxData.ContractID];
	if !ok {
		return errors.New("Not exist contract ID")
	}

	//_, err := os.Stat(f)
	//if os.IsNotExist(err) {
	//	return false
	//}

	return nil
}


func (scs *SmartContractService) Invoke() {

}

func (scs *SmartContractService) keyByValue(ContractPath string) (key string, ok bool) {
	contractName := strings.Replace(ContractPath, "/", "^", -1)
	for k, v := range scs.SmartContractMap {
		if contractName == v.ContractPath {
			key = k
			ok = true
			return key, ok
		}
	}
	return "", false
}
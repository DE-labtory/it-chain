package service

import (
	"errors"
	"it-chain/domain"
	"strings"
	"os"
	"time"
	"context"
	"docker.io/go-docker"
	"io"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"encoding/json"
	"bufio"
	"os/exec"
	"it-chain/common"
	"fmt"
	"github.com/spf13/viper"
)

var logger_s = common.GetLogger("smart_contract_service.go")

type SmartContract struct {
	Name              string
	OriginReposPath   string
	SmartContractPath string
}

type SmartContractServiceImpl struct {
	GithubID              string
	SmartContractHomePath string
	SmartContractMap      map[string]SmartContract
	WorldStateDBPath      string
	WorldStateDBName      string
}

func Init() {

}

func NewSmartContractService(githubID string, smartContractDirPath string) *SmartContractServiceImpl {
	return &SmartContractServiceImpl{
		GithubID:              githubID,
		SmartContractHomePath: smartContractDirPath,
		SmartContractMap:      make(map[string]SmartContract),
	}
}

func (scs *SmartContractServiceImpl) PullAllSmartContracts(errorHandler func(error), completionHandler func()) {
	GOPATH := os.Getenv("GOPATH")
	//go func() {
		repoList, err := domain.GetRepositoryList(scs.GithubID)
		if err != nil {
			errorHandler(errors.New("An error was occurred during getting repository list"))
			return
		}
		for _, repo := range repoList {
			localReposPath := GOPATH + "/src/it-chain" + scs.SmartContractHomePath + "/" + repo.Name

			err = os.MkdirAll(localReposPath, 0755)
			if err != nil {
				errorHandler(errors.New("An error was occurred during making repository path"))
				return
			}

			commits, err := domain.GetReposCommits(repo.FullName)
			if err != nil {
				errorHandler(errors.New("An error was occurred during getting commit logs"))
				return
			}

			for _, commit := range commits {
				err := domain.CloneReposWithName(repo.FullName, localReposPath, commit.Sha)
				if err != nil {
					errorHandler(errors.New("An error was occurred during cloning with name"))
					return
				}

				err = domain.ResetWithSHA(localReposPath+"/"+commit.Sha, commit.Sha)
				if err != nil {
					errorHandler(errors.New("An error was occurred during resetting with SHA"))
					return
				}
			}
		}
		if completionHandler != nil {
			completionHandler()
		}
		return
	//}()

}

func (scs *SmartContractServiceImpl) Deploy(ReposPath string) (string, error) {
	origin_repos_name := strings.Split(ReposPath, "/")[1]
	new_repos_name := strings.Replace(ReposPath, "/", "_", -1)

	_, ok := scs.keyByValue(ReposPath)
	if ok {
		// 버전 업데이트 기능 추가 필요
		return "", errors.New("Already exist smart contract ID")
	}

	repos, err := domain.GetRepos(ReposPath)
	if err != nil {
		return "", errors.New("An error occurred while getting repos!")
	}
	if repos.Message == "Bad credentials" {
		return "", errors.New("Not Exist Repos!")
	}

	err = os.MkdirAll(scs.SmartContractHomePath+"/"+new_repos_name, 0755)
	if err != nil {
		return "", errors.New("An error occurred while make repository's directory!")
	}

	//todo gitpath이미 존재하는지 확인
	err = domain.CloneRepos(ReposPath, scs.SmartContractHomePath+"/"+new_repos_name)
	if err != nil {
		return "", errors.New("An error occurred while cloning repos!")
	}

	common.Log.Println(viper.GetString("smartContract.githubID"))
	_, err = domain.CreateRepos(new_repos_name, viper.GetString("smartContract.githubAccessToken"))
	if err != nil {
		return "", errors.New(err.Error()) //"An error occurred while creating repos!")
	}

	err = domain.ChangeRemote(scs.GithubID+"/"+new_repos_name, scs.SmartContractHomePath+"/"+new_repos_name+"/"+origin_repos_name)
	if err != nil {
		return "", errors.New("An error occurred while cloning repos!")
	}

	// 버전 관리를 위한 파일 추가
	now := time.Now().Format("2006-01-02 15:04:05");
	file, err := os.OpenFile(scs.SmartContractHomePath+"/"+new_repos_name+"/"+origin_repos_name+"/version", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return "", errors.New("An error occurred while creating or opening file!")
	}

	_, err = file.WriteString("Deployed at " + now + "\n")
	if err != nil {
		return "", errors.New("An error occurred while writing file!")
	}
	err = file.Close()
	if err != nil {
		return "", errors.New("An error occurred while closing file!")
	}

	err = domain.CommitAndPush(scs.SmartContractHomePath+"/"+new_repos_name+"/"+origin_repos_name, "It-Chain Smart Contract \""+new_repos_name+"\" Deploy")
	if err != nil {
		return "", errors.New(err.Error())
		//return "", errors.New("An error occurred while committing and pushing!")
	}

	githubResponseCommits, err := domain.GetReposCommits(scs.GithubID + "/" + new_repos_name)
	if err != nil {
		return "", errors.New("An error occurred while getting commit log!")
	}

	reposDirPath := scs.SmartContractHomePath + "/" + new_repos_name + "/" + githubResponseCommits[0].Sha
	err = os.Rename(scs.SmartContractHomePath+"/"+new_repos_name+"/"+origin_repos_name, reposDirPath)
	if err != nil {
		return "", errors.New("An error occurred while renaming directory!")
	}

	scs.SmartContractMap[githubResponseCommits[0].Sha] = SmartContract{new_repos_name, ReposPath, ""}

	return githubResponseCommits[0].Sha, nil
}

/***************************************************
 *	1. smartcontract 검사
 *	2. smartcontract -> sc.tar : 애초에 풀 받을 때 압축해 둘 수 있음
 *	3. go 버전에 맞는 docker image를 Create
 *	4. sc.tar를 docker container로 복사
 *	5. docker container Start
 *	6. docker에서 smartcontract 실행
 ****************************************************/
func (scs *SmartContractServiceImpl) ValidateTransactionsOfBlock(block *domain.Block) (error) {
	// 블럭 유효성 검사 필요?
	if block.TransactionCount <= 0 {
		return errors.New("No tx in block")
	}
	for i := 0; i < block.TransactionCount; i++ {
		scs.ValidateTransaction(block.Transactions[i])
	}
	return nil
}

func (scs *SmartContractServiceImpl) ValidateTransaction(transaction *domain.Transaction) {
	smartContractResponse, err := scs.RunTransactionOnDocker(transaction)
	if err != nil {
		logger_s.Errorln("An error occurred while validating smartcontract!")
		transaction.TransactionStatus = domain.Status_TRANSACTION_UNCONFIRMED
	} else {
		if smartContractResponse.Result == domain.SUCCESS {
			logger_s.Println("Running smartcontract is success")
			transaction.TransactionStatus = domain.Status_TRANSACTION_CONFIRMED
		} else if smartContractResponse.Result == domain.FAIL {
			logger_s.Errorln("An error occurred while validating smartcontract!")
			transaction.TransactionStatus = domain.Status_TRANSACTION_UNCONFIRMED
		}
	}

	transaction.GenerateHash()
}

func (scs *SmartContractServiceImpl) RunTransactionsOfBlock(block domain.Block) (error) {
	if block.TransactionCount <= 0 {
		return errors.New("No tx in block")
	}
	var err error
	for i := 0; i < block.TransactionCount; i++ {
		err = scs.RunTransaction(block.Transactions[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (scs *SmartContractServiceImpl) RunTransaction(transaction *domain.Transaction) (error) {
	if transaction.TxData.Method == domain.Query {
		//smartContractResponse, err := scs.RunTransactionOnDocker(transaction)
		_, err := scs.RunTransactionOnDocker(transaction)
		if err != nil {
			logger_s.Errorln("An error occurred while running smartcontract!")
			return err
		}
	} else if transaction.TxData.Method == domain.Invoke {
		scs.Invoke(transaction)
	}

	return nil
}

func (scs *SmartContractServiceImpl) RunTransactionOnDocker(transaction *domain.Transaction) (*domain.SmartContractResponse, error) {
	GOPATH := os.Getenv("GOPATH")

	/*** Set Transaction Arg ***/
	logger_s.Errorln("validateTransaction start")
	txBytes, err := json.Marshal(transaction)
	if err != nil {
		return nil, errors.New("Tx Marshal Error")
	}

	sc, ok := scs.SmartContractMap[transaction.TxData.ContractID]
	if !ok {
		logger_s.Errorln("Not exist contract ID")
		return nil, errors.New("Not exist contract ID")
	}

	_, err = os.Stat(GOPATH + "/src/it-chain" + sc.SmartContractPath + "/" + transaction.TxData.ContractID)
	if os.IsNotExist(err) {
		logger_s.Errorln("File or Directory Not Exist")
		return nil, errors.New("File or Directory Not Exist")
	}

	/* Docker Code */
	imageName := "docker.io/library/golang:1.9.2-alpine3.6"

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		logger_s.Errorln("An error occurred while creating new Docker Client!")
		return nil, err
	}

	logger_s.Errorln("Pulling image")
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		logger_s.Errorln("An error oeccured while pulling docker image!")
		return nil, err
	}
	io.Copy(os.Stdout, out)

	imageName_splited := strings.Split(imageName, "/")
	image := imageName_splited[len(imageName_splited)-1]

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd: []string{
			"go", "run",
			"/go/src/it-chain" + sc.SmartContractPath + "/" + transaction.TxData.ContractID + "/" + sc.Name + ".go",
			string(txBytes), "/go/src/it-chain" + scs.WorldStateDBPath + "/" + scs.WorldStateDBName,
		},
		//Cmd:          []string{"/go/src/it-chain/smartcontract/sample_smartcontract/" + sc.Name, string(txBytes), GOPATH + "/src" + scs.WorldStateDBPath + "/" + scs.WorldStateDBName},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Binds: []string{GOPATH + "/src/it-chain:/go/src/it-chain"}, // {"/smartcontract/path:/smartcontract", "/worldstatedb/path:/worldstatedb", "/config:/config"},
	}, nil, "")
	if err != nil {
		logger_s.Errorln("An error occurred while creating docker container!")
		return nil, err
	}

	/*** read tar file ***/
	/*
	file, err := ioutil.ReadFile(tarPath)
	if err != nil {
		logger_s.Errorln("An error occurred while reading smartcontract tar file!")
		return nil, err
	}
	wsdb, err := ioutil.ReadFile(tarPath_wsdb)
	if err != nil {
		logger_s.Errorln("An error occurred while reading worldstateDB tar file!")
		return nil, err
	}
	config, err := ioutil.ReadFile(tarPath_config)
	if err != nil {
		logger_s.Errorln("An error occurred while reading config tar file!")
		return nil, err
	}
	*/

	/*** copy file to docker ***/
	/*
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(file), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger_s.Errorln("An error occurred while copying the smartcontract to the container!")
		return nil, err
	}
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(wsdb), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger_s.Errorln("An error occurred while copying the worldstateDB to the container!")
		return nil, err
	}
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(config), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger_s.Errorln("An error occurred while copying the config to the container!")
		return nil, err
	}
	*/

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		logger_s.Errorln("An error occurred while starting the container!")
		return nil, err
	}

	/* get docker output
	----------------------*/
	fmt.Println("============<Docker Output>=============")
	reader, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	if err != nil {
		logger_s.Errorln("An error occurred while getting the output!")
		return nil, err
	}
	defer reader.Close()

	var output = ""
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		output += scanner.Text() + "\n"
	}
	fmt.Println(output)

	smartContractResponse := &domain.SmartContractResponse{}
	err = json.Unmarshal([]byte(output), smartContractResponse)
	if err != nil {
		logger_s.Errorln("An error occurred while unmarshalling the smartContractResponse!")
		return nil, err
	}

	return smartContractResponse, nil
}

func (scs *SmartContractServiceImpl) Query() {
	logger_s.Errorln("Query has to run on docker")
}

func (scs *SmartContractServiceImpl) Invoke(transaction *domain.Transaction) (*domain.SmartContractResponse, error) {
	GOPATH := os.Getenv("GOPATH")

	/*** Set Transaction Arg ***/
	logger_s.Errorln("invoke start")
	txBytes, err := json.Marshal(transaction)
	if err != nil {
		return nil, errors.New("Tx Marshal Error")
	}

	sc, ok := scs.SmartContractMap[transaction.TxData.ContractID]
	if !ok {
		logger_s.Errorln("Not exist contract ID")
		return nil, errors.New("Not exist contract ID")
	}

	_, err = os.Stat(GOPATH + "/src/it-chain" + sc.SmartContractPath + "/" + transaction.TxData.ContractID)
	if os.IsNotExist(err) {
		logger_s.Errorln("File or Directory Not Exist")
		return nil, errors.New("File or Directory Not Exist")
	}

	cmd := exec.Command(
		"go", "run",
		GOPATH + "/src/it-chain" + sc.SmartContractPath + "/" + transaction.TxData.ContractID + "/" + sc.Name + ".go",
		string(txBytes), GOPATH + "/src/it-chain" + scs.WorldStateDBPath + "/" + scs.WorldStateDBName)
	//cmd.Dir = sc.SmartContractPath + "/" + transaction.TxData.ContractID

	output, err := cmd.Output()
	if err != nil {
		logger_s.Errorln("SmartContract running error")
		return nil, err
	}

	fmt.Println("========< Output >========")
	fmt.Println(string(output))

	smartContractResponse := &domain.SmartContractResponse{}
	err = json.Unmarshal(output, smartContractResponse)
	if err != nil {
		logger_s.Errorln("An error occurred while unmarshalling the smartContractResponse!")
		return nil, err
	}
	return smartContractResponse, nil
}

func (scs *SmartContractServiceImpl) keyByValue(OriginReposPath string) (key string, ok bool) {
	contractName := strings.Replace(OriginReposPath, "/", "^", -1)
	for k, v := range scs.SmartContractMap {
		if contractName == v.OriginReposPath {
			key = k
			ok = true
			return key, ok
		}
	}
	return "", false
}

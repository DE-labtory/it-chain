package service

import (
	"errors"
	"it-chain/domain"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"bytes"
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

const (
	TMP_DIR string = "/tmp"
)

var logger_s = common.GetLogger("smart_contract_service.go")

type SmartContract struct {
	Name string
	OriginReposPath string
	SmartContractPath string
}

type SmartContractServiceImpl struct {
	GithubID string
	SmartContractDirPath string
	SmartContractMap map[string]SmartContract
}

func Init() {

}

func NewSmartContractService(githubID string,smartContractDirPath string) *SmartContractServiceImpl{
	return &SmartContractServiceImpl{
		GithubID:githubID,
		SmartContractDirPath:smartContractDirPath,
		SmartContractMap: make(map[string]SmartContract),
	}
}

func (scs *SmartContractServiceImpl) PullAllSmartContracts(authenticatedGit string, errorHandler func(error),
	completionHandler func()) {

	go func() {
		repoList, err := domain.GetRepositoryList(authenticatedGit)
		if err != nil {
			errorHandler(errors.New("An error was occured during getting repository list"))
			return
		}

		for _, repo := range repoList {
			localReposPath := scs.SmartContractDirPath + "/" +
				strings.Replace(repo.FullName, "/", "_", -1)

			err = os.MkdirAll(localReposPath, 0755)
			if err != nil {
				errorHandler(errors.New("An error was occured during making repository path"))
				return
			}

			commits, err := domain.GetReposCommits(repo.FullName)
			if err != nil {
				errorHandler(errors.New("An error was occured during getting commit logs"))
				return
			}

			for _, commit := range commits {
				if commit.Author.Login == authenticatedGit {

					err := domain.CloneReposWithName(repo.FullName, localReposPath, commit.Sha)
					if err != nil {
						errorHandler(errors.New("An error was occured during cloning with name"))
						return
					}

					err = domain.ResetWithSHA(localReposPath + "/" + commit.Sha, commit.Sha)
					if err != nil {
						errorHandler(errors.New("An error was occured during resetting with SHA"))
						return
					}

				}
			}
		}

		completionHandler()
		return
	}()

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
		return "", errors.New("An error occured while getting repos!")
	}
	if repos.Message == "Bad credentials" {
		return "", errors.New("Not Exist Repos!")
	}

	err = os.MkdirAll(scs.SmartContractDirPath + "/" + new_repos_name, 0755)
	if err != nil {
		return "", errors.New("An error occured while make repository's directory!")
	}

	//todo gitpath이미 존재하는지 확인
	err = domain.CloneRepos(ReposPath, scs.SmartContractDirPath + "/" + new_repos_name)
	if err != nil {
		return "", errors.New("An error occured while cloning repos!")
	}

	common.Log.Println(viper.GetString("smartContract.githubID"))
	_, err = domain.CreateRepos(new_repos_name,  viper.GetString("smartContract.githubAccessToken"))
	if err != nil {
		return "", errors.New(err.Error())//"An error occured while creating repos!")
	}

	err = domain.ChangeRemote(scs.GithubID + "/" + new_repos_name, scs.SmartContractDirPath + "/" + new_repos_name + "/" + origin_repos_name)
	if err != nil {
		return "", errors.New("An error occured while cloning repos!")
	}

	// 버전 관리를 위한 파일 추가
	now := time.Now().Format("2006-01-02 15:04:05");
	file, err := os.OpenFile(scs.SmartContractDirPath + "/" + new_repos_name + "/" + origin_repos_name + "/version", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

	err = domain.CommitAndPush(scs.SmartContractDirPath + "/" + new_repos_name + "/" + origin_repos_name, "It-Chain Smart Contract \"" + new_repos_name + "\" Deploy")
	if err != nil {
		return "", errors.New(err.Error())
		//return "", errors.New("An error occured while committing and pushing!")
	}

	githubResponseCommits, err := domain.GetReposCommits(scs.GithubID + "/" + new_repos_name)
	if err != nil {
		return "", errors.New("An error occured while getting commit log!")
	}


	reposDirPath := scs.SmartContractDirPath + "/" + new_repos_name + "/" + githubResponseCommits[0].Sha
	err = os.Rename(scs.SmartContractDirPath + "/" + new_repos_name + "/" + origin_repos_name, reposDirPath)
	if err != nil {
		return "", errors.New("An error occured while renaming directory!")
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
func (scs *SmartContractServiceImpl) Query(transaction domain.Transaction) (error) {
	/*** Set Transaction Arg ***/
	logger_s.Errorln("query start")
	tx_bytes, err := json.Marshal(transaction)
	if err != nil {
		return errors.New("Tx Marshal Error")
	}

	sc, ok := scs.SmartContractMap[transaction.TxData.ContractID];
	if !ok {
		logger_s.Errorln("Not exist contract ID")
		return errors.New("Not exist contract ID")
	}

	_, err = os.Stat(sc.SmartContractPath)
	if os.IsNotExist(err) {
		logger_s.Errorln("File or Directory Not Exist")
		return errors.New("File or Directory Not Exist")
	}

	/*** smartcontract build ***/
	logger_s.Errorln("build start")
	cmd := exec.Command("env", "GOOS=linux", "go", "build", "-o", TMP_DIR + "/" + sc.Name, "./" + sc.Name + ".go")
	cmd.Dir = sc.SmartContractPath + "/" + transaction.TxData.ContractID

	err = cmd.Run()
	if err != nil {
		logger_s.Errorln("SmartContract build error")
		return err
	}
	cmd = exec.Command("chmod", "777", TMP_DIR + "/" + sc.Name)
	cmd.Dir = sc.SmartContractPath + "/" + transaction.TxData.ContractID
	err = cmd.Run()
	if err != nil {
		logger_s.Errorln("Chmod Error")
		return err
	}

	logger_s.Errorln("make tar")

	err = domain.MakeTar(TMP_DIR + "/" + sc.Name, TMP_DIR)
	if err != nil {
		logger_s.Errorln("An error occured while archiving smartcontract file!")
		return err
	}
	err = domain.MakeTar("$GOPATH/src/it-chain/smartcontract/worldstatedb", TMP_DIR)
	if err != nil {
		logger_s.Errorln("An error occured while archiving worldstateDB file!")
		return err
	}

	logger_s.Errorln("exec cmd")

	// tar config file
	cmd = exec.Command("tar", "-cf", TMP_DIR + "/config.tar", "./it-chain/config.yaml")
	cmd.Dir = "../../"
	err = cmd.Run()
	if err != nil {
		logger_s.Errorln("An error occured while archiving config file!")
		return err
	}

	logger_s.Errorln("Pulling image")

	// Docker Code
	imageName := "docker.io/library/golang:1.9.2-alpine3.6"
	tarPath := TMP_DIR + "/" + sc.Name + ".tar"
	tarPath_wsdb := TMP_DIR + "/worldstatedb.tar"
	tarPath_config := TMP_DIR + "/config.tar"

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		logger_s.Errorln("An error occured while creating new Docker Client!")
		return err
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		logger_s.Errorln("An error oeccured while pulling docker image!")
		return err
	}
	io.Copy(os.Stdout, out)


	imageName_splited := strings.Split(imageName, "/")
	image := imageName_splited[len(imageName_splited)-1]

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd: []string{"/go/src/" + sc.Name, string(tx_bytes)},
		Tty: true,
		AttachStdout: true,
		AttachStderr: true,
	}, nil, nil, "")
	if err != nil {
		logger_s.Errorln("An error occured while creating docker container!")
		return err
	}

	/*** read tar file ***/
	file, err := ioutil.ReadFile(tarPath)
	if err != nil {
		logger_s.Errorln("An error occured while reading smartcontract tar file!")
		return err
	}
	wsdb, err := ioutil.ReadFile(tarPath_wsdb)
	if err != nil {
		logger_s.Errorln("An error occured while reading worldstateDB tar file!")
		return err
	}
	config, err := ioutil.ReadFile(tarPath_config)
	if err != nil {
		logger_s.Errorln("An error occured while reading config tar file!")
		return err
	}

	/*** copy file to docker ***/
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(file), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger_s.Errorln("An error occured while copying the smartcontract to the container!")
		return err
	}
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(wsdb), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger_s.Errorln("An error occured while copying the worldstateDB to the container!")
		return err
	}
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(config), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger_s.Errorln("An error occured while copying the config to the container!")
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		logger_s.Errorln("An error occured while starting the container!")
		return err
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
		logger_s.Errorln("An error occured while getting the output!")
		return err
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
	fmt.Println("----Marshaled Output----")
	fmt.Println(smartContractResponse)

	if smartContractResponse.Result == domain.SUCCESS {
		logger_s.Println("Running smartcontract is success")
		transaction.GenerateHash()
		// real running smartcontract
	} else if smartContractResponse.Result == domain.FAIL {
		logger_s.Errorln("An error occured while running smartcontract!")
	}

	return nil
}

func (scs *SmartContractServiceImpl) Invoke() {

}

func (scs *SmartContractServiceImpl) ValidateTransactionsInBlock(block *domain.Block) error {
	// 블럭 유효성 검사 필요?
	if block.TransactionCount <= 0 {
		return errors.New("No tx in block")
	}
	var err error
	for i := 0; i < block.TransactionCount; i++ {
		err = scs.ValidateTransaction(block.Transactions[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (scs *SmartContractServiceImpl) ValidateTransaction(transaction *domain.Transaction) (error) {
	/*** Set Transaction Arg ***/
	logger_s.Errorln("validateTransaction start")
	txBytes, err := json.Marshal(transaction)
	if err != nil {
		return errors.New("Tx Marshal Error")
	}

	sc, ok := scs.SmartContractMap[transaction.TxData.ContractID]
	if !ok {
		logger_s.Errorln("Not exist contract ID")
		return errors.New("Not exist contract ID")
	}

	_, err = os.Stat(sc.SmartContractPath)
	if os.IsNotExist(err) {
		logger_s.Errorln("File or Directory Not Exist")
		return errors.New("File or Directory Not Exist")
	}

	/*** smartcontract build ***/
	logger_s.Errorln("build start")
	fmt.Println(sc.SmartContractPath + "/" + transaction.TxData.ContractID + "/" + sc.Name + ".go")
	cmd := exec.Command("env", "GOOS=linux", "go", "build", "-o", TMP_DIR + "/" + sc.Name, sc.SmartContractPath + "/" + transaction.TxData.ContractID + "/" + sc.Name + ".go")
	cmd.Dir = sc.SmartContractPath + "/" + transaction.TxData.ContractID

	err = cmd.Run()
	if err != nil {
		logger_s.Errorln("SmartContract build error")
		return err
	}
	cmd = exec.Command("chmod", "777", TMP_DIR + "/" + sc.Name)
	cmd.Dir = sc.SmartContractPath + "/" + transaction.TxData.ContractID
	err = cmd.Run()
	if err != nil {
		logger_s.Errorln("Chmod Error")
		return err
	}

	logger_s.Errorln("make tar")

	err = domain.MakeTar(TMP_DIR + "/" + sc.Name, TMP_DIR)
	if err != nil {
		logger_s.Errorln("An error occured while archiving smartcontract file!")
		return err
	}
	err = domain.MakeTar("$GOPATH/src/it-chain/smartcontract/worldstatedb", TMP_DIR)
	if err != nil {
		logger_s.Errorln("An error occured while archiving worldstateDB file!")
		return err
	}

	logger_s.Errorln("exec cmd")

	// tar config file
	cmd = exec.Command("tar", "-cf", TMP_DIR + "/config.tar", "./it-chain/config.yaml")
	cmd.Dir = "../../"
	err = cmd.Run()
	if err != nil {
		logger_s.Errorln("An error occured while archiving config file!")
		return err
	}

	logger_s.Errorln("Pulling image")

	// Docker Code
	imageName := "docker.io/library/golang:1.9.2-alpine3.6"
	tarPath := TMP_DIR + "/" + sc.Name + ".tar"
	tarPath_wsdb := TMP_DIR + "/worldstatedb.tar"
	tarPath_config := TMP_DIR + "/config.tar"

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		logger_s.Errorln("An error occured while creating new Docker Client!")
		return err
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		logger_s.Errorln("An error oeccured while pulling docker image!")
		return err
	}
	io.Copy(os.Stdout, out)


	imageName_splited := strings.Split(imageName, "/")
	image := imageName_splited[len(imageName_splited)-1]

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd: []string{"/go/src/" + sc.Name, string(txBytes)},
		Tty: true,
		AttachStdout: true,
		AttachStderr: true,
	}, nil, nil, "")
	if err != nil {
		logger_s.Errorln("An error occured while creating docker container!")
		return err
	}

	/*** read tar file ***/
	file, err := ioutil.ReadFile(tarPath)
	if err != nil {
		logger_s.Errorln("An error occured while reading smartcontract tar file!")
		return err
	}
	wsdb, err := ioutil.ReadFile(tarPath_wsdb)
	if err != nil {
		logger_s.Errorln("An error occured while reading worldstateDB tar file!")
		return err
	}
	config, err := ioutil.ReadFile(tarPath_config)
	if err != nil {
		logger_s.Errorln("An error occured while reading config tar file!")
		return err
	}

	/*** copy file to docker ***/
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(file), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger_s.Errorln("An error occured while copying the smartcontract to the container!")
		return err
	}
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(wsdb), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger_s.Errorln("An error occured while copying the worldstateDB to the container!")
		return err
	}
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(config), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger_s.Errorln("An error occured while copying the config to the container!")
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		logger_s.Errorln("An error occured while starting the container!")
		return err
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
		logger_s.Errorln("An error occured while getting the output!")
		return err
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

	if smartContractResponse.Result == domain.SUCCESS {
		logger_s.Println("Running smartcontract is success")
		transaction.GenerateHash()
	} else if smartContractResponse.Result == domain.FAIL {
		logger_s.Errorln("An error occured while validating smartcontract!")
		return errors.New("An error occured while validating smartcontract!")
	}
	return nil
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
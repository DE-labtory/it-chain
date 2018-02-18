package smartcontract

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
)

const (
	GITHUB_TOKEN string = "31d8f4c1bfc6906806b9a77803087b5b671fac2d"
	TMP_DIR string = "/tmp"
)
var logger = common.GetLogger("smart_contract_service.go")

// temp
const AUTHENTICATED_GIT = "emperorhan"

type SmartContract struct {
	Name string
	OriginReposPath string
	SmartContractPath string
}

type SmartContractService struct {
	GithubID string
	SmartContractDirPath string
	SmartContractMap map[string]SmartContract
}

func Init() {

}

func (scs *SmartContractService) pullAllSmartContracts() (error) {

	repoList, err := GetRepositoryList(AUTHENTICATED_GIT)
	if err != nil {
		return errors.New("An error was occured during getting repository list")
	}

	for _, repo := range repoList {
		localReposPath := scs.SmartContractDirPath + "/" +
			strings.Replace(repo.FullName, "/", "_", -1)

		err = os.MkdirAll(localReposPath, 0755)
		if err != nil {
			return errors.New("An error was occured during making repository path")
		}

		commits, err := GetReposCommits(repo.FullName)
		if err != nil {
			return errors.New("An error was occured during getting commit logs")
		}

		for _, commit := range commits {
			if commit.Author.Login == AUTHENTICATED_GIT {

				err := CloneReposWithName(repo.FullName, localReposPath, commit.Sha)
				if err != nil {
					return errors.New("An error was occured during cloning with name")
				}

				err = ResetWithSHA(localReposPath + "/" + commit.Sha, commit.Sha)
				if err != nil {
					return errors.New("An error was occured during resetting with SHA")
				}

			}
		}
	}

	return nil
}

func (scs *SmartContractService) Deploy(ReposPath string) (string, error) {
	origin_repos_name := strings.Split(ReposPath, "/")[1]
	new_repos_name := strings.Replace(ReposPath, "/", "_", -1)

	_, ok := scs.keyByValue(ReposPath)
	if ok {
		// 버전 업데이트 기능 추가 필요
		return "", errors.New("Already exist smart contract ID")
	}

	repos, err := GetRepos(ReposPath)
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

	err = CloneRepos(ReposPath, scs.SmartContractDirPath + "/" + new_repos_name)
	if err != nil {
		return "", errors.New("An error occured while cloning repos!")
	}

	_, err = CreateRepos(new_repos_name, GITHUB_TOKEN)
	if err != nil {
		return "", errors.New(err.Error())//"An error occured while creating repos!")
	}

	err = ChangeRemote(scs.GithubID + "/" + new_repos_name, scs.SmartContractDirPath + "/" + new_repos_name + "/" + origin_repos_name)
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

	err = CommitAndPush(scs.SmartContractDirPath + "/" + new_repos_name + "/" + origin_repos_name, "It-Chain Smart Contract \"" + new_repos_name + "\" Deploy")
	if err != nil {
		return "", errors.New(err.Error())
		//return "", errors.New("An error occured while committing and pushing!")
	}

	githubResponseCommits, err := GetReposCommits(scs.GithubID + "/" + new_repos_name)
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
func (scs *SmartContractService) Query(transaction domain.Transaction) (error) {
	/*** Set Transaction Arg ***/
	tx_bytes, err := json.Marshal(transaction)
	if err != nil {
		return errors.New("Tx Marshal Error")
	}

	sc, ok := scs.SmartContractMap[transaction.TxData.ContractID];
	if !ok {
		logger.Errorln("Not exist contract ID")
		return errors.New("Not exist contract ID")
	}

	_, err = os.Stat(sc.SmartContractPath)
	if os.IsNotExist(err) {
		logger.Errorln("File or Directory Not Exist")
		return errors.New("File or Directory Not Exist")
	}

	/*** smartcontract build ***/
	cmd := exec.Command("env", "GOOS=linux", "go", "build", "-o", TMP_DIR + "/" + sc.Name, "./" + sc.Name + ".go")
	cmd.Dir = sc.SmartContractPath + "/" + transaction.TxData.ContractID
	err = cmd.Run()
	if err != nil {
		logger.Errorln("SmartContract build error")
		return err
	}
	cmd = exec.Command("chmod", "777", TMP_DIR + "/" + sc.Name)
	cmd.Dir = sc.SmartContractPath + "/" + transaction.TxData.ContractID
	err = cmd.Run()
	if err != nil {
		logger.Errorln("Chmod Error")
		return err
	}

	err = MakeTar(TMP_DIR + "/" + sc.Name, TMP_DIR)
	if err != nil {
		logger.Errorln("An error occured while archiving smartcontract file!")
		return err
	}
	err = MakeTar("$GOPATH/src/it-chain/smartcontract/worldstatedb", TMP_DIR)
	if err != nil {
		logger.Errorln("An error occured while archiving worldstateDB file!")
		return err
	}

	// tar config file
	cmd = exec.Command("tar", "-cf", TMP_DIR + "/config.tar", "./it-chain/config.yaml")
	cmd.Dir = "../../"
	err = cmd.Run()
	if err != nil {
		logger.Errorln("An error occured while archiving config file!")
		return err
	}

	// Docker Code
	imageName := "docker.io/library/golang:1.9.2-alpine3.6"
	tarPath := TMP_DIR + "/" + sc.Name + ".tar"
	tarPath_wsdb := TMP_DIR + "/worldstatedb.tar"
	tarPath_config := TMP_DIR + "/config.tar"

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		logger.Errorln("An error occured while creating new Docker Client!")
		return err
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		logger.Errorln("An error oeccured while pulling docker image!")
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
		logger.Errorln("An error occured while creating docker container!")
		return err
	}

	/*** read tar file ***/
	file, err := ioutil.ReadFile(tarPath)
	if err != nil {
		logger.Errorln("An error occured while reading smartcontract tar file!")
		return err
	}
	wsdb, err := ioutil.ReadFile(tarPath_wsdb)
	if err != nil {
		logger.Errorln("An error occured while reading worldstateDB tar file!")
		return err
	}
	config, err := ioutil.ReadFile(tarPath_config)
	if err != nil {
		logger.Errorln("An error occured while reading config tar file!")
		return err
	}

	/*** copy file to docker ***/
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(file), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger.Errorln("An error occured while copying the smartcontract to the container!")
		return err
	}
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(wsdb), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger.Errorln("An error occured while copying the worldstateDB to the container!")
		return err
	}
	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(config), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		logger.Errorln("An error occured while copying the config to the container!")
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		logger.Errorln("An error occured while starting the container!")
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
		logger.Errorln("An error occured while getting the output!")
		return err
	}
	defer reader.Close()

	var output = ""
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		output += scanner.Text() + "\n"
	}
	fmt.Println(output)

	smartContractResponse := &SmartContractResponse{}
	err = json.Unmarshal([]byte(output), smartContractResponse)
	fmt.Println("----Marshaled Output----")
	fmt.Println(smartContractResponse)

	if smartContractResponse.Result == SUCCESS {
		logger.Println("Running smartcontract is success")

		// tx hash reset
		// real running smartcontract
	} else if smartContractResponse.Result == FAIL {
		logger.Errorln("An error occured while running smartcontract!")
	}

	return nil
}


func (scs *SmartContractService) Invoke() {

}

func (scs *SmartContractService) keyByValue(OriginReposPath string) (key string, ok bool) {
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
package smartcontract

import (
	"errors"
	"it-chain/service/blockchain"
	"strings"
	"fmt"
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
)

const (
	GITHUB_TOKEN string = "31d8f4c1bfc6906806b9a77803087b5b671fac2d"

)

// temp
const AUTHENTICATED_GIT = "emperorhan"

type SmartContract struct {
	ReposName string
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
func (scs *SmartContractService) Query(transaction blockchain.Transaction) (error) {
	fmt.Println("func Query Start")
	tmpDir := "/tmp"
	sc, ok := scs.SmartContractMap[transaction.TxData.ContractID];
	if !ok {
		fmt.Println("Not exist contract ID")
		return errors.New("Not exist contract ID")
	}

	_, err := os.Stat(sc.SmartContractPath)
	if os.IsNotExist(err) {
		fmt.Println("File or Directory Not Exist")
		return errors.New("File or Directory Not Exist")
	}

	err = MakeTar(sc.SmartContractPath + "/" + transaction.TxData.ContractID, tmpDir)
	if err != nil {
		return errors.New("An error occured while archiving file!")
	}
	fmt.Println("Passed MakeTar Smart Contract")

	err = MakeTar("/Users/hackurity/go/src/it-chain/smartcontract/worldstatedb", tmpDir)
	if err != nil {
		return errors.New("An error occured while archiving file!")
	}
	fmt.Println("Passed MakeTar World State DB")

	fmt.Println("======== sc =======")
	fmt.Println(sc)

	//PullAndCopyAndRunDocker("docker.io/library/golang:rc-alpine", tmpDir+"/"+transaction.TxData.ContractID+".tar")

	// Docker Code
	imageName := "docker.io/library/golang:rc-alpine"
	tarPath := tmpDir+"/"+transaction.TxData.ContractID+".tar"
	tarPath_wsdb := tmpDir+"/worldstatedb.tar"

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
	fmt.Println("Passed ImagePull")

	imageName_splited := strings.Split(imageName, "/")
	image := imageName_splited[len(imageName_splited)-1]

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		//Cmd: []string{"/bin/bash"},
		Tty: true,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	fmt.Println("Passed ContainerCreate")

	file, err := ioutil.ReadFile(tarPath)
	if err != nil {
		fmt.Print(err)
	}

	wsdb, err := ioutil.ReadFile(tarPath_wsdb)
	if err != nil {
		fmt.Print(err)
	}

	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(file), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Passed CopyToContainer Go File")

	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(wsdb), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Passed CopyToContainer World State DB")

	fmt.Println("============================")
	fmt.Println("resp.ID : " + resp.ID)
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Passed ContainerStart")


	/* Set Transaction Arg
	------------------------*/
	tx_bytes, err := json.Marshal(transaction)
	if err != nil {
		return errors.New("Tx Marshal Error")
	}
	fmt.Println("Passed Marshal Tx")


	/* go build in docker
	----------------------*/
	exec_build, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
		Cmd: []string{"go", "build", "-o", "/go/src/" + transaction.TxData.ContractID + "/" + sc.ReposName, "/go/src/" + transaction.TxData.ContractID + "/" + sc.ReposName + ".go"},
		User: "root",
	})	// go build -o /go/src/abc/sample1 /go/src/abc/sample1.go
	if err != nil {
		panic(err)
	}

	err = cli.ContainerExecStart(ctx, exec_build.ID, types.ExecStartCheck{
		Detach: true,
		Tty:    true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Passed go build in docker")

	fmt.Println("------------ tx_byte ------------")
	fmt.Println(string(tx_bytes))

	/* go run in docker
	--------------------*/
	exec_run, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
		Cmd: []string{"/go/src/" + transaction.TxData.ContractID + "/" + sc.ReposName, string(tx_bytes)},
		User: "root",
	})
	if err != nil {
		panic(err)
	}

	err = cli.ContainerExecStart(ctx, exec_run.ID, types.ExecStartCheck{
		Detach: true,
		Tty:    true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Passed go run in docker")

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
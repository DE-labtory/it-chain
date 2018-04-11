package git

import (
	"fmt"
	"os/user"

	"strings"

	"github.com/it-chain/it-chain-Engine/smartcontract/domain/smartContract"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var tmp = "./.tmp"

type GitApi struct {
	sshAuth *ssh.PublicKeys
}

func NewGitApi() GitApi {

	currentUser, err := user.Current()

	if err != nil {
		panic(fmt.Sprintf("fail to init GitApi [%s]", err.Error()))
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "")

	if err != nil {
		panic(fmt.Sprintf("fail to init GitApi [%s]", err.Error()))
	}

	return GitApi{
		sshAuth: sshAuth,
	}
}

//todo SSH ENV로 ssh key 불러오기
func (g GitApi) Clone(gitUrl string) (*smartContract.SmartContract, error) {

	r, err := git.PlainClone(tmp, false, &git.CloneOptions{
		URL:               gitUrl,
		Auth:              g.sshAuth,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	committer, err := r.CommitObjects()
	lastCommit, err := committer.Next()
	commitHash := lastCommit.Hash.String()

	if err != nil {
		return nil, err
	}

	//todo os separator
	name := getNameFromGitUrl(gitUrl)
	sc := smartContract.NewSmartContract(name, gitUrl, tmp+"/"+name, commitHash)

	return sc, nil
}

//
func getNameFromGitUrl(gitUrl string) string {
	parsed := strings.Split(gitUrl, "/")

	if len(parsed) == 0 {
		return ""
	}

	name := parsed[len(parsed)-1]

	return name
}

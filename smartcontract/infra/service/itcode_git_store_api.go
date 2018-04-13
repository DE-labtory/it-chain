package service

import (
	"fmt"
	"os/user"
	"strings"

	"os"

	"github.com/it-chain/it-chain-Engine/smartcontract/domain/itcode"
	"github.com/pkg/errors"
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

//get itcode from outside
//todo SSH ENV로 ssh key 불러오기
func (g GitApi) Clone(gitUrl string) (*itcode.ItCode, error) {

	r, err := git.PlainClone(tmp, false, &git.CloneOptions{
		URL:               gitUrl,
		Auth:              g.sshAuth,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	head, err := r.Head()

	if err != nil {
		return nil, err
	}

	lastHeadCommit, err := r.CommitObject(head.Hash())
	commitHash := lastHeadCommit.Hash.String()

	if err != nil {
		return nil, err
	}

	//todo os separator
	name := getNameFromGitUrl(gitUrl)
	sc := itcode.NewItCode(name, gitUrl, tmp+"/"+name, commitHash)

	return sc, nil
}

//push code to auth repo
func (g GitApi) Push(itCode itcode.ItCode) error {
	itCodePath := itCode.Path

	if !dirExists(itCodePath) {
		return errors.New(fmt.Sprintf("Invalid itCode Path [%s]", itCodePath))
	}

	return nil
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

func dirExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}

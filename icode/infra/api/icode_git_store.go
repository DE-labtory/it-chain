package api

import (
	"strings"

	"context"

	"os/user"

	"errors"
	"fmt"

	"os"

	"github.com/google/go-github/github"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/icode"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type ICodeGitStoreApi struct {
	AuthClient  *github.Client
	AuthGitUser *github.User
	authUserID  string
	authUserPW  string
	name        string
}

func NewICodeGitStoreApi(authUserID string, authUserPW string) (*ICodeGitStoreApi, error) {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(authUserID),
		Password: strings.TrimSpace(authUserPW),
	}
	client := github.NewClient(tp.Client())
	ctx := context.Background()
	gitUser, _, err := client.Users.Get(ctx, "")

	if err != nil {
		return nil, err
	}

	confSshPath := conf.GetConfiguration().Icode.SshPath
	if confSshPath == "default" {
		currentUser, err := user.Current()
		if err != nil {
			return nil, err
		}
		confSshPath = currentUser.HomeDir + "/.ssh/id_rsa"

	}

	return &ICodeGitStoreApi{
		AuthClient:  client,
		AuthGitUser: gitUser,
		authUserID:  authUserID,
		authUserPW:  authUserPW,
		name:        gitUser.GetLogin(),
	}, nil
}

func (gApi *ICodeGitStoreApi) Clone(baseSavePath string, repositoryUrl string, sshPath string) (*icode.Meta, error) {
	name := getNameFromGitUrl(repositoryUrl)

	if name == "" {
		return nil, errors.New(fmt.Sprintf("Invalid url name [%s]", repositoryUrl))
	}

	//check file already exist
	if _, err := os.Stat(baseSavePath + "/" + name); err == nil {
		// if iCode already exist, remove that
		err = os.RemoveAll(baseSavePath + "/" + name)
		if err != nil {
			return nil, err
		}
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", sshPath, "")
	if err != nil {
		return nil, err
	}

	r, err := git.PlainClone(baseSavePath+"/"+name, false, &git.CloneOptions{
		URL:               repositoryUrl,
		Auth:              sshAuth,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		return nil, err
	}

	head, err := r.Head()

	if err != nil {
		return nil, err
	}

	lastHeadCommit, err := r.CommitObject(head.Hash())
	commitHash := lastHeadCommit.Hash.String()

	if err != nil {
		return nil, err
	}

	metaData := icode.NewMeta(name, repositoryUrl, baseSavePath+"/"+name, commitHash)
	return metaData, nil
}

//todo async push
func (gApi *ICodeGitStoreApi) Push(meta icode.Meta) error {
	iCodePath := meta.Path

	if _, err := os.Stat(iCodePath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Invalid iCodeMeta Path [%s]", iCodePath))
	}

	err := gApi.createRepository(meta.RepositoryName)

	if err != nil {
		return err
	}

	err = changeRemote(iCodePath, gApi.AuthGitUser.GetHTMLURL()+"/"+meta.RepositoryName)

	if err != nil {
		return err
	}

	au := &http.BasicAuth{Username: gApi.authUserID, Password: gApi.authUserPW}

	r, err := git.PlainOpen(iCodePath)

	if err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{
		RemoteName: git.DefaultRemoteName,
		Auth:       au,
	})

	if err != nil {
		return err
	}

	return nil

}

func (gApi *ICodeGitStoreApi) createRepository(name string) error {
	ctx := context.Background()
	_, _, err := gApi.AuthClient.Repositories.Create(ctx, "", &github.Repository{
		Name:    &name,
		Private: &[]bool{false}[0],
	})
	return err
}

func getNameFromGitUrl(gitUrl string) string {
	parsed := strings.Split(gitUrl, "/")

	if len(parsed) == 0 {
		return ""
	}

	name := strings.Split(parsed[len(parsed)-1], ".")

	if len(name) == 0 {
		return ""
	}

	return name[0]
}

func changeRemote(path string, changeUrl string) error {
	r, err := git.PlainOpen(path)

	if err != nil {
		return err
	}

	err = r.DeleteRemote(git.DefaultRemoteName)

	if err != nil {
		return err
	}

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{changeUrl},
	})

	if err != nil {
		return err
	}

	return nil
}

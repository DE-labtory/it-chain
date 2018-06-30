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
	authClient  *github.Client
	AuthGitUser *github.User
	authUserID  string
	authUserPW  string
	name        string
	sshAuth     *ssh.PublicKeys
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

	sshAuth, err := ssh.NewPublicKeysFromFile("git", confSshPath, "")
	if err != nil {
		return nil, err
	}

	return &ICodeGitStoreApi{
		authClient:  client,
		AuthGitUser: gitUser,
		authUserID:  authUserID,
		authUserPW:  authUserPW,
		name:        gitUser.GetLogin(),
		sshAuth:     sshAuth,
	}, nil
}

func (gApi *ICodeGitStoreApi) Clone(repositoryUrl string) (*icode.Meta, error) {
	name := getNameFromGitUrl(repositoryUrl)

	if name == "" {
		return nil, errors.New(fmt.Sprintf("Invalid url name [%s]", repositoryUrl))
	}

	r, err := git.PlainClone(conf.GetConfiguration().Icode.ICodeSavePath+"/"+name, false, &git.CloneOptions{
		URL:               repositoryUrl,
		Auth:              gApi.sshAuth,
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

	metaData := icode.NewMeta(name, repositoryUrl, conf.GetConfiguration().Icode.ICodeSavePath+"/"+name, commitHash)
	return metaData, nil
}

//todo async push
func (gApi *ICodeGitStoreApi) Push(meta icode.Meta) error {
	iCodePath := meta.Path

	if _, err := os.Stat(iCodePath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Invalid iCodeMeta Path [%s]", iCodePath))
	}

	err := changeRemote(iCodePath, gApi.AuthGitUser.GetHTMLURL()+"/"+meta.RepositoryName)

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

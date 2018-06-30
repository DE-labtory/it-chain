package api

import (
	"strings"

	"context"

	"os/user"

	"github.com/google/go-github/github"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/icode"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type ICodeGitStoreApi struct {
	authClient *github.Client
	authUserID string
	authUserPW string
	name       string
	sshAuth    *ssh.PublicKeys
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
		authClient: client,
		authUserID: authUserID,
		authUserPW: authUserPW,
		name:       gitUser.GetLogin(),
		sshAuth:    sshAuth,
	}, nil
}

func (*ICodeGitStoreApi) Clone(repositoryUrl string) (*icode.Meta, error) {
	panic("implement me")
}

func (*ICodeGitStoreApi) Push(meta icode.Meta) error {
	panic("implement me")
}

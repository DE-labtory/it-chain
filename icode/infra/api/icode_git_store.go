package api

import (
	"github.com/google/go-github/github"
	"github.com/it-chain/it-chain-Engine/icode"
)

type ICodeGitStoreApi struct {
	authClient *github.Client
}

func NewICodeGitStoreApi(authUserID string, authUserPW string) (*ICodeGitStoreApi, error) {
	panic("implement please")
}

func (*ICodeGitStoreApi) Clone(repositoryUrl string) (*icode.Meta, error) {
	panic("implement me")
}

func (*ICodeGitStoreApi) Push(meta icode.Meta) error {
	panic("implement me")
}

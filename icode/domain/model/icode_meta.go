package model

import "github.com/rs/xid"

type Version struct {
}

type ICodeID string

type ICodeMeta struct {
	ID             ICodeID
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
}

func NewICodeMeta(repositoryName string, gitUrl string, path string, commitHash string) *ICodeMeta {

	return &ICodeMeta{
		ID:             ICodeID(xid.New().String()),
		RepositoryName: repositoryName,
		CommitHash:     commitHash,
		GitUrl:         gitUrl,
		Path:           path,
	}
}

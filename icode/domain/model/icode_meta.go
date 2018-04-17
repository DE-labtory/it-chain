package model

import (
	"encoding/json"

	"github.com/rs/xid"
)

type Version struct {
}

type ICodeID string

func (i ICodeID) ToString() string {
	return string(i)
}

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

func (i ICodeMeta) Serialize() ([]byte, error) {
	b, err := json.Marshal(i)

	if err != nil {
		return nil, err
	}

	return b, nil
}

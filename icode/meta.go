package icode

import (
	"errors"
	"fmt"

	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type Version struct {
}

type ID = string

type Meta struct {
	ICodeID        ID
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
}

func NewMeta(repositoryName string, gitUrl string, path string, commitHash string) *Meta {
	return &Meta{
		ICodeID:        xid.New().String(),
		RepositoryName: repositoryName,
		CommitHash:     commitHash,
		GitUrl:         gitUrl,
		Path:           path,
	}
}

func (m Meta) GetID() string {
	return m.ICodeID
}

func (m *Meta) On(event midgard.Event) error {
	switch v := event.(type) {

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

type MetaRepository interface {
	Save(meta Meta) error
	Remove(id ID) error
	FindById(id ID) (*Meta, error)
	FindAll() ([]*Meta, error)
}

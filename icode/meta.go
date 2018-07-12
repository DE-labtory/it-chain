package icode

import (
	"errors"
	"fmt"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
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
	createEvent := MetaCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   xid.New().String(),
			Type: "meta.created",
		},
		RepositoryName: repositoryName,
		GitUrl:         gitUrl,
		Path:           path,
		CommitHash:     commitHash,
	}
	eventstore.Save(createEvent.GetID(), createEvent)
	m := Meta{}
	m.On(&createEvent)
	return &m
}

func DeleteMeta(metaId ID) error {
	return eventstore.Save(metaId, MetaDeletedEvent{
		EventModel: midgard.EventModel{
			ID:   metaId,
			Type: "meta.deleted",
		},
	})
}

func (m Meta) GetID() string {
	return m.ICodeID
}

func (m *Meta) On(event midgard.Event) error {
	switch v := event.(type) {
	case *MetaCreatedEvent:
		m.CommitHash = v.CommitHash
		m.Version = v.Version
		m.Path = v.Path
		m.GitUrl = v.GitUrl
		m.ICodeID = v.ID
		m.RepositoryName = v.RepositoryName
	case *MetaDeletedEvent:
		m.ICodeID = ""
		m.RepositoryName = ""
		m.GitUrl = ""
		m.Path = ""
		m.CommitHash = ""
		m.Version = Version{}
	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

type MetaRepository interface {
	Save(meta Meta) error
	Remove(id ID) error
	FindById(id ID) (*Meta, error)
	FindByGitURL(url string) (*Meta, error)
	FindAll() ([]*Meta, error)
}

type ReadOnlyMetaRepository interface {
	FindById(id ID) (*Meta, error)
	FindByGitURL(url string) (*Meta, error)
	FindAll() ([]*Meta, error)
}

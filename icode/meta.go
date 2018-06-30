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
	return createEvent.GetMeta()
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

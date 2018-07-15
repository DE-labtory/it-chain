package icode

import (
	"errors"
	"fmt"

	"time"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type Version struct {
}

type ID = string
type MetaStatus = int

const (
	READY = iota
	UNDEPLOYED
	DEPLOYED
	DEPLOY_FAIL
)

type Meta struct {
	ICodeID        ID
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
	Status         MetaStatus
}

func NewMeta(id string, repositoryName string, gitUrl string, path string, commitHash string) *Meta {
	createEvent := MetaCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   id,
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

func ChangeMetaStatus(metaId ID, status MetaStatus) error {
	return eventstore.Save(metaId, MetaStatusChangeEvent{
		EventModel: midgard.EventModel{
			ID:   metaId,
			Type: "meta.statusChange",
			Time: time.Now(),
		},
		Status: DEPLOYED,
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
		m.Status = READY
	case *MetaDeletedEvent:
		m.Status = UNDEPLOYED
	case *MetaStatusChangeEvent:
		m.Status = v.Status
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

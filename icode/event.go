package icode

import "github.com/it-chain/midgard"

type MetaCreatedEvent struct {
	midgard.EventModel
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
}

func (m MetaCreatedEvent) GetMeta() *Meta {
	return &Meta{
		ICodeID:        m.ID,
		RepositoryName: m.RepositoryName,
		GitUrl:         m.GitUrl,
		Path:           m.Path,
		CommitHash:     m.CommitHash,
		Version:        m.Version,
	}
}

type MetaDeletedEvent struct {
	midgard.EventModel
}

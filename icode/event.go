package icode

import "github.com/it-chain/midgard"

//type : meta.created
type MetaCreatedEvent struct {
	midgard.EventModel
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
}

//type : meta.deleted
type MetaDeletedEvent struct {
	midgard.EventModel
}

type MetaStatusChangeEvent struct {
	midgard.EventModel
	Status MetaStatus
}

/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package icode

import (
	"errors"
	"fmt"

	"time"

	"github.com/it-chain/engine/core/eventstore"
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
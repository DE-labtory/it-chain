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

package api_gateway

import (
	"errors"
	"log"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/leveldb-wrapper"
)

type ICodeQueryApi struct {
	metaRepository ICodeMetaRepository
}

func NewICodeQueryApi(repository ICodeMetaRepository) ICodeQueryApi {
	return ICodeQueryApi{
		metaRepository: repository,
	}
}

type ICodeMetaRepository interface {
	FindAllMeta() ([]ivm.Meta, error)
	FindMetaByUrl(url string) (ivm.Meta, error)
	FindMetaById(id ivm.ID) (ivm.Meta, error)
	Save(meta ivm.Meta) error
	Remove(id ivm.ID) error
}

type LevelDbMetaRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewLevelDbMetaRepository(path string) LevelDbMetaRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return LevelDbMetaRepository{
		leveldb: db,
	}
}

func (l *LevelDbMetaRepository) FindAllMeta() ([]ivm.Meta, error) {
	iter := l.leveldb.GetIteratorWithPrefix([]byte(""))
	metaList := []ivm.Meta{}
	for iter.Next() {
		val := iter.Value()
		meta := &ivm.Meta{}
		err := common.Deserialize(val, meta)
		if err != nil {
			return nil, err
		}
		metaList = append(metaList, *meta)
	}
	return metaList, nil
}

func (l *LevelDbMetaRepository) FindMetaByUrl(url string) (ivm.Meta, error) {
	allMetaList, err := l.FindAllMeta()
	if err != nil {
		return ivm.Meta{}, err
	}

	for _, meta := range allMetaList {
		if meta.GitUrl == url {
			return meta, nil
		}
	}

	return ivm.Meta{}, nil
}

func (l *LevelDbMetaRepository) FindMetaById(id ivm.ID) (ivm.Meta, error) {

	metaByte, err := l.leveldb.Get([]byte(id))
	if err != nil {
		return ivm.Meta{}, err
	}

	if len(metaByte) == 0 {
		return ivm.Meta{}, nil
	}

	meta := &ivm.Meta{}

	err = common.Deserialize(metaByte, meta)

	if err != nil {
		return ivm.Meta{}, err
	}

	return *meta, nil
}

func (l *LevelDbMetaRepository) Save(meta ivm.Meta) error {

	if meta.ICodeID == "" {
		return errors.New("meta is empty")
	}

	b, err := common.Serialize(meta)
	if err != nil {
		return err
	}

	err = l.leveldb.Put([]byte(meta.ICodeID), b, true)
	if err != nil {
		return err
	}

	return nil
}

func (l *LevelDbMetaRepository) Remove(id ivm.ID) error {
	return l.leveldb.Delete([]byte(id), true)
}

type ICodeEventHandler struct {
	metaRepository ICodeMetaRepository
}

func NewIcodeEventHandler(repository ICodeMetaRepository) ICodeEventHandler {
	return ICodeEventHandler{
		metaRepository: repository,
	}
}

func (i ICodeEventHandler) HandleMetaCreatedEvent(metaCreatedEvent event.MetaCreated) {

	meta := ivm.Meta{
		ICodeID:        metaCreatedEvent.ICodeID,
		RepositoryName: metaCreatedEvent.RepositoryName,
		GitUrl:         metaCreatedEvent.GitUrl,
		Path:           metaCreatedEvent.Path,
		CommitHash:     metaCreatedEvent.CommitHash,
		Version:        metaCreatedEvent.Version,
	}

	err := i.metaRepository.Save(meta)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (i ICodeEventHandler) HandleMetaDeletedEvent(metaDeleted event.MetaDeleted) {

	err := i.metaRepository.Remove(metaDeleted.ICodeID)

	if err != nil {
		log.Fatal(err.Error())
	}
}

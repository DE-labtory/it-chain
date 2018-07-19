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
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/engine/common"
	"errors"
	"log"
	"fmt"
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
	FindAllMeta() ([]icode.Meta, error)
	FindMetaByUrl(url string) (icode.Meta, error)
	FindMetaById(id icode.ID) (icode.Meta, error)
	Save(meta icode.Meta) error
	Remove(id icode.ID) error
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

func (l *LevelDbMetaRepository) FindAllMeta() ([]icode.Meta, error) {
	iter := l.leveldb.GetIteratorWithPrefix([]byte(""))
	metaList := []icode.Meta{}
	for iter.Next() {
		val := iter.Value()
		meta := &icode.Meta{}
		err := common.Deserialize(val, meta)
		if err != nil {
			return nil, err
		}
		metaList = append(metaList, *meta)
	}
	return metaList, nil
}

func (l *LevelDbMetaRepository) FindMetaByUrl(url string) (icode.Meta, error) {
	allMetaList, err := l.FindAllMeta()
	if err != nil {
		return icode.Meta{}, err
	}

	for _, meta := range allMetaList {
		if meta.GitUrl == url {
			return meta, nil
		}
	}

	return icode.Meta{}, nil
}

func (l *LevelDbMetaRepository) FindMetaById(id icode.ID) (icode.Meta, error) {

	metaByte, err := l.leveldb.Get([]byte(id))
	if err != nil {
		return icode.Meta{}, err
	}

	if len(metaByte) == 0 {
		return icode.Meta{}, nil
	}

	meta := &icode.Meta{}

	err = common.Deserialize(metaByte,meta)

	if err != nil {
		return icode.Meta{}, err
	}

	return *meta, nil
}

func (l *LevelDbMetaRepository) Save(meta icode.Meta) error {
	if meta.GetID() == "" {
		return errors.New("meta is empty")
	}
	b, err := common.Serialize(meta)
	if err != nil {
		return err
	}

	err = l.leveldb.Put([]byte(meta.GetID()), b, true)
	if err != nil {
		return err
	}
	return nil
}

func (l *LevelDbMetaRepository) Remove(id icode.ID) error {
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

func (i ICodeEventHandler) HandleMetaCreatedEvent(event icode.MetaCreatedEvent) {
	meta := event.GetMeta()
	err := i.metaRepository.Save(meta)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (i ICodeEventHandler) HandleMetaDeletedEvent(event icode.MetaDeletedEvent) {
	err := i.metaRepository.Remove(event.GetID())

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (i ICodeEventHandler) HandleMetaStatusChangeEvent(event icode.MetaStatusChangeEvent) {
	meta, err := i.metaRepository.FindMetaById(event.GetID())
	if err != nil {
		log.Fatal(err.Error())
	}
	if meta.GetID() ==""{
		log.Fatal(fmt.Sprintf("no icode id : [%s] in handleMetaStatusChangeEvent(api_gateway/icode_query_api",event.GetID()))
	}
	meta.Status = event.Status
	err = i.metaRepository.Save(meta)
	if err != nil {
		log.Fatal(err.Error())
	}
}

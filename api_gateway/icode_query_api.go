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
	iCodeRepository ICodeRepository
}

func NewICodeQueryApi(repository ICodeRepository) ICodeQueryApi {
	return ICodeQueryApi{
		iCodeRepository: repository,
	}
}

type ICodeRepository interface {
	FindAllMeta() ([]ivm.ICode, error)
	FindMetaByUrl(url string) (ivm.ICode, error)
	FindMetaById(id ivm.ID) (ivm.ICode, error)
	Save(icode ivm.ICode) error
	Remove(id ivm.ID) error
	Close()
}

type LevelDbICodeRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewLevelDbMetaRepository(path string) *LevelDbICodeRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &LevelDbICodeRepository{
		leveldb: db,
	}
}

func (l *LevelDbICodeRepository) FindAllMeta() ([]ivm.ICode, error) {
	iter := l.leveldb.GetIteratorWithPrefix([]byte(""))
	metaList := []ivm.ICode{}
	for iter.Next() {
		val := iter.Value()
		icode := &ivm.ICode{}
		err := common.Deserialize(val, icode)
		if err != nil {
			return nil, err
		}
		metaList = append(metaList, *icode)
	}
	return metaList, nil
}

func (l *LevelDbICodeRepository) FindMetaByUrl(url string) (ivm.ICode, error) {
	allMetaList, err := l.FindAllMeta()
	if err != nil {
		return ivm.ICode{}, err
	}

	for _, icode := range allMetaList {
		if icode.GitUrl == url {
			return icode, nil
		}
	}

	return ivm.ICode{}, nil
}

func (l *LevelDbICodeRepository) FindMetaById(id ivm.ID) (ivm.ICode, error) {

	metaByte, err := l.leveldb.Get([]byte(id))
	if err != nil {
		return ivm.ICode{}, err
	}

	if len(metaByte) == 0 {
		return ivm.ICode{}, nil
	}

	icode := &ivm.ICode{}

	err = common.Deserialize(metaByte, icode)

	if err != nil {
		return ivm.ICode{}, err
	}

	return *icode, nil
}

func (l *LevelDbICodeRepository) Save(icode ivm.ICode) error {

	if icode.ID == "" {
		return errors.New("icode is empty")
	}

	b, err := common.Serialize(icode)
	if err != nil {
		return err
	}

	err = l.leveldb.Put([]byte(icode.ID), b, true)
	if err != nil {
		return err
	}

	return nil
}

func (l *LevelDbICodeRepository) Remove(id ivm.ID) error {
	return l.leveldb.Delete([]byte(id), true)
}

func (l *LevelDbICodeRepository) Close() {
	l.leveldb.Close()
}

type ICodeEventHandler struct {
	iCodeRepository ICodeRepository
}

func NewIcodeEventHandler(repository ICodeRepository) *ICodeEventHandler {
	return &ICodeEventHandler{
		iCodeRepository: repository,
	}
}

func (i ICodeEventHandler) HandleMetaCreatedEvent(icodeCreatedEvent event.ICodeCreated) {

	icode := ivm.ICode{
		ID:             icodeCreatedEvent.ID,
		RepositoryName: icodeCreatedEvent.RepositoryName,
		GitUrl:         icodeCreatedEvent.GitUrl,
		Path:           icodeCreatedEvent.Path,
		CommitHash:     icodeCreatedEvent.CommitHash,
		Version:        icodeCreatedEvent.Version,
	}

	err := i.iCodeRepository.Save(icode)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (i ICodeEventHandler) HandleMetaDeletedEvent(iCodeDeleted event.ICodeDeleted) {

	err := i.iCodeRepository.Remove(iCodeDeleted.ICodeID)

	if err != nil {
		log.Fatal(err.Error())
	}
}

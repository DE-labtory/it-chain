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
	"os"
	"testing"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/ivm"
	"github.com/stretchr/testify/assert"
)

func TestLevelDbMetaRepository_Save(t *testing.T) {
	// setting
	dbPath := "./.test"
	repo := NewLevelDbMetaRepository(dbPath)
	defer func() {
		repo.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	// given
	tests := map[string]struct {
		Input  ivm.Meta
		Output error
	}{
		"success": {
			Input: ivm.Meta{
				ICodeID:        "1",
				RepositoryName: "name",
				GitUrl:         "url",
				Path:           "path",
				CommitHash:     "hash",
				Version:        ivm.Version{},
			},
			Output: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running '%s' test, caseName: %s", t.Name(), testName)
		//given
		outputError := repo.Save(test.Input)
		//then
		assert.Equal(t, test.Output, outputError, "error in save")
		//check
		b, err := repo.leveldb.Get([]byte(test.Input.ICodeID))
		assert.NoError(t, err)
		checkMeta := &ivm.Meta{}
		assert.NoError(t, err, "error in checking process, leveldb get")
		err = common.Deserialize(b, checkMeta)
		assert.NoError(t, err, "error in checking process, deserialize")
		assert.Equal(t, test.Input, *checkMeta)
	}
}

func TestLevelDbMetaRepository_FindAllMeta(t *testing.T) {
	// setting
	dbPath := "./.test"
	repo := NewLevelDbMetaRepository(dbPath)
	defer func() {
		repo.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	tests := map[string]struct {
		SettingData []ivm.Meta
		Output      error
	}{
		"success": {
			SettingData: []ivm.Meta{{
				ICodeID:        "1",
				RepositoryName: "a",
				GitUrl:         "a",
				Path:           "a",
				CommitHash:     "a",
			}, {
				ICodeID:        "2",
				RepositoryName: "b",
				GitUrl:         "b",
				Path:           "b",
				CommitHash:     "b",
			}},
			Output: nil,
		},
	}
	for testName, test := range tests {
		t.Logf("Running '%s' test, caseName: %s", t.Name(), testName)
		//given
		for _, data := range test.SettingData {
			err := repo.Save(data)
			assert.NoError(t, err, "error in setting data")
		}
		resultDatas, err := repo.FindAllMeta()
		assert.NoError(t, err, "error in find all")
		//then
		assert.Equal(t, len(test.SettingData), len(resultDatas))
	}

}

func TestLevelDbMetaRepository_FindMetaById(t *testing.T) {
	// setting
	dbPath := "./.test"
	repo := NewLevelDbMetaRepository(dbPath)
	defer func() {
		repo.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	setData := ivm.Meta{
		ICodeID:        "123",
		RepositoryName: "a",
		GitUrl:         "a",
		Path:           "a",
		CommitHash:     "a",
	}

	err := repo.Save(setData)
	assert.NoError(t, err, "error while setting data")

	//setting map
	tests := map[string]struct {
		Input       ivm.ID
		Output      ivm.Meta
		OutputError error
	}{
		"success": {
			Input:       "123",
			Output:      setData,
			OutputError: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running '%s' test, caseName: %s", t.Name(), testName)

		//given
		meta, err := repo.FindMetaById("123")

		//then
		assert.Equal(t, test.Output, meta)
		assert.Equal(t, test.OutputError, err)
	}
}

func TestLevelDbMetaRepository_FindMetaByUrl(t *testing.T) {
	// setting
	dbPath := "./.test"
	repo := NewLevelDbMetaRepository(dbPath)
	defer func() {
		repo.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	setData := ivm.Meta{
		ICodeID:        "123",
		RepositoryName: "a",
		GitUrl:         "gitUrl",
		Path:           "a",
		CommitHash:     "a",
	}

	err := repo.Save(setData)
	assert.NoError(t, err, "error while setting data")

	//setting map
	tests := map[string]struct {
		Input       string
		Output      ivm.Meta
		OutputError error
	}{
		"success": {
			Input:       "gitUrl",
			Output:      setData,
			OutputError: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running '%s' test, caseName: %s", t.Name(), testName)

		//given
		meta, err := repo.FindMetaByUrl("gitUrl")

		//then
		assert.Equal(t, test.Output, meta)
		assert.Equal(t, test.OutputError, err)
	}
}

// todo test fail 함
//func TestICodeEventHandler_HandleMetaCreatedEvent(t *testing.T) {
//
//	git, client, tearDown := setICodeQueryApi(t)
//
//	defer tearDown()
//
//	tests := map[string]struct {
//		Input         ivm.MetaCreatedEvent
//		OutputError   error
//		ExpectDataNum int
//	}{
//		"success": {
//			Input: ivm.MetaCreatedEvent{
//				EventModel: midgard.EventModel{
//					ID:   "1",
//					Type: "meta.created",
//					Time: time.Now(),
//				},
//				RepositoryName: "a",
//				GitUrl:         "b",
//				Path:           "c",
//				CommitHash:     "d",
//			},
//			OutputError:   nil,
//			ExpectDataNum: 1,
//		},
//	}
//
//	for testName, test := range tests {
//		t.Logf("Running '%s' test, caseName: %s", t.Name(), testName)
//		//given
//		err := client.Publish("Event", "transaction.created", test.Input)
//		time.Sleep(3 * time.Second)
//
//		//then
//		assert.Equal(t, test.OutputError, err, "err in compare err")
//
//		//check
//		metas, err := git.metaRepository.FindAllMeta()
//		assert.NoError(t, err, "err in check")
//		assert.Equal(t, test.ExpectDataNum, len(metas), "not equal in check dataNum")
//
//	}
//}

func TestICodeEventHandler_HandleMetaStatusChangeEvent(t *testing.T) {
	//todo impl like TestICodeEventHandler_HandleMetaCreatedEvent
}

func TestICodeEventHandler_HandleMetaDeletedEvent(t *testing.T) {
	//todo impl like TestICodeEventHandler_HandleMetaCreatedEvent
}

func setICodeQueryApi(t *testing.T) (ICodeQueryApi, *pubsub.TopicPublisher, func()) {

	dbPath := "./.test"
	client := pubsub.NewTopicSubscriber("", "Event")
	publisher := pubsub.NewTopicPublisher("", "Event")

	repo := NewLevelDbMetaRepository(dbPath)

	metaQueryApi := ICodeQueryApi{metaRepository: &repo}
	metaEventListener := &ICodeEventHandler{metaRepository: &repo}

	err := client.SubscribeTopic("meta.*", metaEventListener)
	assert.NoError(t, err)

	return metaQueryApi, &publisher, func() {
		os.RemoveAll(dbPath)
		client.Close()
	}
}

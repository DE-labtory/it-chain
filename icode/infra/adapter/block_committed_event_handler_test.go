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

package adapter_test

import (
	"errors"
	"testing"

	"os"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/engine/icode/api"
	"github.com/it-chain/engine/icode/infra/adapter"
	"github.com/it-chain/engine/icode/infra/git"
	"github.com/it-chain/engine/icode/infra/tesseract"
	"github.com/stretchr/testify/assert"
)

func TestBlockCommittedEventHandler_HandleBlockCommittedEventHandler_When_Event_Occur(t *testing.T) {

}

func TestBlockCommittedEventHandler_HandleBlockCommittedEventHandler(t *testing.T) {

	//given
	handler, containerService, tearDown := setUp(t)
	defer tearDown()

	testBlock := event.BlockCommitted{
		TxList: []event.Tx{
			event.Tx{
				ICodeID:  "1",
				Function: "initA",
				Args:     []string{},
			},
		},
	}

	//when
	handler.HandleBlockCommittedEventHandler(testBlock)

	//then
	// success case
	result, err := containerService.ExecuteRequest(icode.Request{
		ICodeID:  "1",
		Function: "getA",
		Type:     "query",
		Args:     []string{},
	})

	assert.NoError(t, err)
	assert.Equal(t, result.Data["A"], "0")
	assert.Equal(t, result.Err, "")
}

// setup handler and start container
func setUp(t *testing.T) (*adapter.BlockCommittedEventHandler, *tesseract.ContainerService, func()) {
	GOPATH := os.Getenv("GOPATH")

	if GOPATH == "" {
		t.Fatal(errors.New("need go path"))
		return nil, nil, func() {}
	}

	// git generate
	storeApi := git.NewRepositoryService()
	containerService := tesseract.NewContainerService()
	eventService := common.NewEventService("", "Event")
	icodeApi := api.NewICodeApi(containerService, storeApi, eventService)

	meta := icode.Meta{
		ICodeID:        "1",
		RepositoryName: "test icode",
		Path:           GOPATH + "/src/github.com/it-chain/engine/icode/mock/",
		GitUrl:         "github.com/mock",
	}

	err := containerService.StartContainer(meta)
	assert.NoError(t, err)

	blockCommittedEventHandler := adapter.NewBlockCommittedEventHandler(icodeApi)

	return blockCommittedEventHandler, containerService, func() {
		containerService.StopContainer(meta.ICodeID)
	}
}

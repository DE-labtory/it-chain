/*
 * Copyright 2018 DE-labtory
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

package api

import (
	"testing"

	"github.com/DE-labtory/engine/ivm"
	"github.com/magiconair/properties/assert"
)

func Test_createMetaCreatedEvent(t *testing.T) {

	icode := ivm.ICode{
		ID:             "1",
		GitUrl:         "https://github.com/DE-labtory/engine",
		RepositoryName: "jun",
		CommitHash:     "hduh48183",
		Path:           "1",
	}

	event := createMetaCreatedEvent(icode)

	assert.Equal(t, event.ID, icode.ID)
	assert.Equal(t, event.CommitHash, icode.CommitHash)
	assert.Equal(t, event.RepositoryName, icode.RepositoryName)
	assert.Equal(t, event.Path, icode.Path)
}

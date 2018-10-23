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

package ipfs_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/it-chain/engine/storage/infra/ipfs"
	"github.com/stretchr/testify/assert"
)

func TestClient_Upload(t *testing.T) {
	client := ipfs.NewClient("")

	fileID, err := client.UploadFile([]byte("hello"))
	assert.NoError(t, err)
	fmt.Println(fileID)
}

func TestClient_DownLoad(t *testing.T) {

	data := []byte("hello")
	client := ipfs.NewClient("")

	fileID, err := client.UploadFile(data)
	assert.NoError(t, err)

	b, err := client.DownLoadFile(fileID.ID)
	defer os.RemoveAll(ipfs.TMP_FOLDER)
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}

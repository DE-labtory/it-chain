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

package ipfs

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/DE-labtory/it-chain/common"
	"github.com/DE-labtory/it-chain/storage"
	shell "github.com/ipfs/go-ipfs-api"
)

const TMP_FOLDER = "./ipfs-tmp/"

type Client struct {
	sh *shell.Shell
}

func NewClient(ipfsAddress string) *Client {

	defaultAddress := "localhost:5001"
	if ipfsAddress == "" {
		ipfsAddress = defaultAddress
	}

	sh := shell.NewShell(ipfsAddress)

	return &Client{
		sh: sh,
	}
}

func (c *Client) UploadFile(b []byte) (storage.FileID, error) {
	cid, err := c.sh.Add(bytes.NewReader(b))
	if err != nil {
		return storage.NewFileID(""), err
	}

	return storage.NewFileID(cid), nil
}

func (c *Client) DownLoadFile(id string) ([]byte, error) {

	err := common.CreateDirIfMissing(TMP_FOLDER)
	if err != nil {
		return []byte{}, err
	}

	path := TMP_FOLDER + id
	defer os.RemoveAll(path)

	err = c.sh.Get(id, path)
	if err != nil {
		return []byte{}, err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

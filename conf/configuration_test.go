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

package conf

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasValid(t *testing.T) {

	isModeFalse := HasValidMode(instance)
	assert.False(t, isModeFalse)

}

func TestGetConfiguration(t *testing.T) {

	path := os.Getenv("GOPATH") + "/src/github.com/it-chain/engine/conf"
	confFilename := "/config-test.yaml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	// yaml이라 탭말고 스페이스로 놔둬주세요 ㅠㅠ자동변환 하는 ide 쓰시는분들 유의 부탁드립니다(vs code 등)
	ConfJson := []byte(`
  engine:
    logpath: log/it-chain.log
    keypath: .it-chain/
    mode: solo
    amqp: amqp://guest:guest@localhost:5672/
    bootstrapnodeaddress: 127.0.0.1:5555
  txpool:
    timeoutms: 1000
    maxtransactionbyte: 1024
  consensus:
    batchtime: 3
    maxtransactions: 100
  blockchain:
    genesisconfpath: ./Genesis.conf
  peer:
    leaderelection: RAFT
  icode:
    repositorypath: empty
  grpcgateway:
    address: 127.0.0.1
    port: "13579"
  apigateway:
    address: 127.0.0.1
    port: "4444"
    `)

	err := ioutil.WriteFile(path+confFilename, ConfJson, os.ModePerm)
	assert.NoError(t, err)

	SetConfigPath(path + confFilename)
	config := GetConfiguration()
	assert.Equal(t, config.GrpcGateway.Address, "127.0.0.1")
	assert.Equal(t, config.Engine.Mode, "solo")

	defer os.Remove(path + confFilename)

}

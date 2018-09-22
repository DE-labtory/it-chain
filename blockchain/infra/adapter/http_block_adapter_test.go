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
	"net/http"
	"testing"

	"time"

	"log"

	"encoding/json"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/stretchr/testify/assert"
)

func TestHttpBlockAdapter_GetLastBlock(t *testing.T) {

	lastBlock := blockchain.DefaultBlock{
		Seal:   []byte("seal"),
		Height: 10,
	}

	go createServer(lastBlock, "8080")

	time.Sleep(4)

	hAdapter := adapter.HttpBlockAdapter{}

	peer := blockchain.Peer{
		ApiGatewayAddress: "127.0.0.1:8080",
	}

	retrievedBlock, err := hAdapter.GetLastBlockFromPeer(peer)

	assert.Equal(t, lastBlock, retrievedBlock)
	assert.NoError(t, err)
}

func TestHttpBlockAdapter_GetLastBlock_Fail_NoLastBlock(t *testing.T) {

	go createServer(blockchain.DefaultBlock{}, "8081")

	time.Sleep(4)

	hAdapter := adapter.HttpBlockAdapter{}

	peer := blockchain.Peer{
		ApiGatewayAddress: "127.0.0.1:8081",
	}

	_, err := hAdapter.GetLastBlockFromPeer(peer)

	assert.Error(t, err)
}

func TestHttpBlockAdapter_GetLastBlock_Fail_NoServer(t *testing.T) {

	hAdapter := adapter.HttpBlockAdapter{}

	peer := blockchain.Peer{
		ApiGatewayAddress: "127.0.0.1:8082",
	}

	_, err := hAdapter.GetLastBlockFromPeer(peer)

	assert.Error(t, err)
}

func TestHttpBlockAdapter_GetBlockByHeight(t *testing.T) {
	_5thBlock := blockchain.DefaultBlock{
		Seal:   []byte("seal"),
		Height: 5,
	}

	go createServer(_5thBlock, "8083")

	time.Sleep(6)

	hAdapter := adapter.HttpBlockAdapter{}

	peer := blockchain.Peer{
		ApiGatewayAddress: "127.0.0.1:8083",
	}

	retrievedBlock, err := hAdapter.GetBlockByHeightFromPeer(_5thBlock.Height, peer)

	assert.Equal(t, _5thBlock, retrievedBlock)
	assert.NoError(t, err)

}

func TestHttpBlockAdapter_GetBlockByHeight_Fail_NoBlockWithTheHeight(t *testing.T) {
	EmptyBlock := blockchain.DefaultBlock{}

	go createServer(EmptyBlock, "8084")

	time.Sleep(6)

	hAdapter := adapter.HttpBlockAdapter{}

	peer := blockchain.Peer{
		ApiGatewayAddress: "127.0.0.1:8084",
	}

	_, err := hAdapter.GetBlockByHeightFromPeer(5, peer)

	//fmt.Println(retrievedBlock)
	assert.Error(t, err)

}

func TestHttpBlockAdapter_GetBlockByHeight_Fail_NoServer(t *testing.T) {

	hAdapter := adapter.HttpBlockAdapter{}

	peer := blockchain.Peer{
		ApiGatewayAddress: "127.0.0.1:8085",
	}

	_, err := hAdapter.GetBlockByHeightFromPeer(5, peer)

	assert.Error(t, err)

}

func createServer(block blockchain.DefaultBlock, port string) {

	b, _ := json.Marshal(block)

	handler := func(w http.ResponseWriter, req *http.Request) {

		if heightStr := req.URL.Query().Get("height"); heightStr != "" {
			w.Write(b)
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/blocks", handler)

	log.Fatal(http.ListenAndServe(":"+port, mux))

}

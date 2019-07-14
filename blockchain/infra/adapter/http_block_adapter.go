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

package adapter

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/DE-labtory/engine/blockchain"
	"gopkg.in/resty.v1"
)

var ErrGetBlockFromPeer = errors.New("Error when getting block from other peer")

type HttpBlockAdapter struct {
}

func NewHttpBlockAdapter() *HttpBlockAdapter {
	return &HttpBlockAdapter{}
}

func (a HttpBlockAdapter) GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.DefaultBlock, error) {

	resp, err := resty.R().
		SetQueryString("height=-1").
		SetHeader("Content-Type", "application/json").
		Get("http://" + peer.ApiGatewayAddress + "/blocks")
	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	block := blockchain.DefaultBlock{}
	if err := json.Unmarshal(resp.Body(), &block); err != nil {
		return blockchain.DefaultBlock{}, err
	}

	if block.IsEmpty() {
		return blockchain.DefaultBlock{}, ErrGetBlockFromPeer
	}

	return block, nil
}

func (a HttpBlockAdapter) GetBlockByHeightFromPeer(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error) {

	resp, err := resty.R().
		SetQueryParams(map[string]string{
			"height": strconv.FormatUint(height, 10),
		}).
		SetHeader("Content-Type", "application/json").
		Get("http://" + peer.ApiGatewayAddress + "/blocks")
	if err != nil {
		return blockchain.DefaultBlock{}, err

	}

	block := blockchain.DefaultBlock{}
	if err := json.Unmarshal(resp.Body(), &block); err != nil {
		return blockchain.DefaultBlock{}, err
	}

	if block.IsEmpty() {
		return blockchain.DefaultBlock{}, ErrGetBlockFromPeer
	}

	return block, nil
}

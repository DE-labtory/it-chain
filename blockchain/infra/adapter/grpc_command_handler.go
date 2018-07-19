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

package adapter

import (
	"encoding/json"

	"github.com/it-chain/engine/blockchain"
)

type SyncBlockApi interface {
	SyncedCheck(block blockchain.Block) error
}

type SyncCheckGrpcCommandService interface {
	SyncCheckResponse(block blockchain.Block) error
	ResponseBlock(peerId blockchain.PeerId, block blockchain.Block) error
}

type GrpcCommandHandler struct {
	blockApi           SyncBlockApi
	blockQueryApi      blockchain.BlockQueryApi
	grpcCommandService SyncCheckGrpcCommandService
}

func NewGrpcCommandHandler(blockApi SyncBlockApi, blockQueryService blockchain.BlockQueryApi, grpcCommandService SyncCheckGrpcCommandService) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		blockApi:           blockApi,
		blockQueryApi:      blockQueryService,
		grpcCommandService: grpcCommandService,
	}
}

func (g *GrpcCommandHandler) HandleGrpcCommand(command blockchain.GrpcReceiveCommand) error {
	switch command.Protocol {
	case "SyncCheckRequestProtocol":
		//TODO: 상대방의 SyncCheck를 위해서 자신의 last block을 보내준다.
		block, err := g.blockQueryApi.GetLastBlock()
		if err != nil {
			return ErrGetLastBlock
		}

		err = g.grpcCommandService.SyncCheckResponse(block)
		if err != nil {
			return ErrSyncCheckResponse
		}
		break

	case "SyncCheckResponseProtocol":
		//TODO: 상대방의 last block을 받아서 SyncCheck를 시작한다.
		break

	case "BlockRequestProtocol":
		var height uint64
		err := json.Unmarshal(command.Body, &height)
		if err != nil {
			return ErrBlockInfoDeliver
		}

		block, err := g.blockQueryApi.GetBlockByHeight(height)
		if err != nil {
			return ErrGetBlock
		}

		err = g.grpcCommandService.ResponseBlock(command.FromPeer.PeerId, block)
		if err != nil {
			return ErrResponseBlock
		}
		break

	case "BlockResponseProtocol":
		//TODO: Construct 과정에서 block을 받는다.
		break
	}

	return nil
}

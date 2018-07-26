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
	blockQueryService  blockchain.BlockQueryService
	grpcCommandService SyncCheckGrpcCommandService
}

func NewGrpcCommandHandler(blockApi SyncBlockApi, blockQueryService blockchain.BlockQueryService, grpcCommandService SyncCheckGrpcCommandService) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		blockApi:           blockApi,
		blockQueryService:  blockQueryService,
		grpcCommandService: grpcCommandService,
	}
}

func (g *GrpcCommandHandler) HandleGrpcCommand(command blockchain.GrpcReceiveCommand) error {

	return nil
}

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
	"testing"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/engine/txpool/infra/adapter"
	"github.com/it-chain/engine/txpool/infra/mem"
	"github.com/it-chain/engine/txpool/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestBlockService_ProposeBlock(t *testing.T) {
	client := rpc.NewClient("")
	server := rpc.NewServer("")

	err := server.Register("block.propose", func(command command.ProposeBlock) (struct{}, rpc.Error) {

		return struct{}{}, rpc.Error{}
	})

	assert.NoError(t, err)

	txpoolRepository := mem.NewTransactionRepository()
	txpoolRepository.Save(txpool.Transaction{ID: "tx1"})
	txpoolRepository.Save(txpool.Transaction{ID: "tx2"})

	peerQueryService := &mock.PeerQueryService{}

	peerQueryService.GetPeerListFunc = func() ([]p2p.Peer, error) {
		return nil, nil
	}

	peerQueryService.GetLeaderFunc = func() (p2p.Leader, error) {
		return p2p.Leader{LeaderId: struct{ Id string }{Id: "1"}}, nil
	}

	myself := command.MyPeer{IpAddress: "", PeerId: "1"}

	transactions, _ := txpoolRepository.FindAll()
	assert.Equal(t, 2, len(transactions))

	blockService := adapter.NewBlockProposalService(client, txpoolRepository, "solo", peerQueryService, myself)
	err = blockService.ProposeBlock()
	assert.NoError(t, err)

	blockService = adapter.NewBlockProposalService(client, txpoolRepository, "pbft", peerQueryService, myself)
	err = blockService.ProposeBlock()
	assert.NoError(t, err)

	myself = command.MyPeer{IpAddress: "", PeerId: "2"}

	blockService = adapter.NewBlockProposalService(client, txpoolRepository, "pbft", peerQueryService, myself)
	err = blockService.ProposeBlock()
	assert.NoError(t, err)

	transactions, _ = txpoolRepository.FindAll()
	assert.Equal(t, 0, len(transactions))
}

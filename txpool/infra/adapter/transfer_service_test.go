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

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/engine/txpool/infra/adapter"
	"github.com/it-chain/engine/txpool/infra/mem"
	"github.com/it-chain/engine/txpool/test/mock"
	"github.com/magiconair/properties/assert"
)

func TestGrpcCommandService_SendLeaderTransactions(t *testing.T) {
	tests := map[string]struct {
		input struct {
			transactions []txpool.Transaction
			leader       txpool.Leader
		}
		err error
	}{
		"success": {
			input: struct {
				transactions []txpool.Transaction
				leader       txpool.Leader
			}{
				transactions: []txpool.Transaction{{ID: txpool.TransactionId("zf")}},
				leader:       txpool.Leader{LeaderId: txpool.LeaderId{Id: "zf2"}},
			},
			err: nil,
		},
		"transaction empty test": {
			input: struct {
				transactions []txpool.Transaction
				leader       txpool.Leader
			}{
				transactions: []txpool.Transaction{},
				leader:       txpool.Leader{LeaderId: txpool.LeaderId{Id: "zf2"}},
			},
			err: adapter.ErrTxEmpty,
		},
	}

	publisher := func(topic string, data interface{}) (err error) {
		txList := &[]*txpool.Transaction{}
		deliverCommand := data.(command.DeliverGrpc)

		common.Deserialize(deliverCommand.Body, txList)

		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, 1, len(*txList))

		return nil
	}

	txpoolRepository := mem.NewTransactionRepository()
	txpoolRepository.Save(txpool.Transaction{ID: "zf"})

	peerQueryService := &mock.PeerQueryService{}

	peerQueryService.GetPeerListFunc = func() ([]p2p.Peer, error) {
		a := []p2p.Peer{}

		a = append(a, p2p.Peer{IpAddress: "12345", PeerId: struct{ Id string }{Id: "zf"}})
		a = append(a, p2p.Peer{IpAddress: "54321", PeerId: struct{ Id string }{Id: "zf2"}})
		return a, nil
	}

	peerQueryService.GetLeaderFunc = func() (p2p.Leader, error) {
		return p2p.Leader{LeaderId: struct{ Id string }{Id: "zf2"}}, nil
	}
	myself := command.MyPeer{IpAddress: "", PeerId: "zf"}

	transferService := adapter.NewTransferService(publisher, txpoolRepository, "pbft", peerQueryService, myself)
	t.Logf("running test case %s", "success")
	err := transferService.SendLeaderTransactions()
	assert.Equal(t, tests["success"].err, err)

	txpoolRepository = mem.NewTransactionRepository()
	txpoolRepository.Save(txpool.Transaction{})

	transferService = adapter.NewTransferService(publisher, txpoolRepository, "pbft", peerQueryService, myself)
	t.Logf("running test case %s", "transaction empty test")
	err = transferService.SendLeaderTransactions()
	assert.Equal(t, tests["transaction empty test"].err, err)
}

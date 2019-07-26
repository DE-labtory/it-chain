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
	"github.com/DE-labtory/iLogger"
	"github.com/DE-labtory/it-chain/blockchain"
	"github.com/DE-labtory/it-chain/common/event"
)

type SynchronizeApi interface {
	HandleNetworkJoined(peer []blockchain.Peer) error
}

type NetworkEventHandler struct {
	SyncApi SynchronizeApi
}

func NewNetworkEventHandler(syncApi SynchronizeApi) *NetworkEventHandler {

	return &NetworkEventHandler{
		SyncApi: syncApi,
	}
}

func (n *NetworkEventHandler) HandleNetworkJoinedEvent(networkJoindEvent event.NetworkJoined) {
	iLogger.Infof(nil, "[Blockchain] Network Joined")
	if err := n.SyncApi.HandleNetworkJoined(createPeerListFromNetworkJoinedEvent(networkJoindEvent)); err != nil {
		iLogger.Errorf(nil, "[Blockchain] Fail to handle network joined event - Err: [%s]", err.Error())
	}
}

func createPeerListFromNetworkJoinedEvent(networkJoindEvent event.NetworkJoined) []blockchain.Peer {
	peerList := make([]blockchain.Peer, 0)

	for _, c := range networkJoindEvent.Connections {
		peerList = append(peerList, blockchain.Peer{
			Id:                c.ConnectionID,
			ApiGatewayAddress: c.ApiGatewayAddress,
		})
	}

	return peerList
}

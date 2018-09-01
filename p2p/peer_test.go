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

package p2p_test

import (
	"testing"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/p2p"
	"github.com/magiconair/properties/assert"
)

func TestPeer_GetID(t *testing.T) {
	peer := &p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "1",
		},
		IpAddress: "1",
	}

	assert.Equal(t, peer.GetID(), "1")
}

func TestPeer_Serialize(t *testing.T) {
	peer := &p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "1",
		},
		IpAddress: "1",
	}

	byte, _ := common.Serialize(peer)
	b, _ := peer.Serialize()

	assert.Equal(t, b, byte)
}

func TestDeserialize(t *testing.T) {
	peer := &p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "1",
		},
	}

	targetPeer := &p2p.Peer{}

	byte, _ := peer.Serialize()

	assert.Equal(t, p2p.Deserialize(byte, targetPeer), nil)
	assert.Equal(t, targetPeer, &p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "1",
		},
	})
}

func TestPeerId_ToString(t *testing.T) {
	peerId := &p2p.PeerId{
		Id: "1",
	}

	assert.Equal(t, peerId.ToString(), "1")
}

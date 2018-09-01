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

package p2p

import (
	"testing"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/magiconair/properties/assert"
)

func TestNewPLTable(t *testing.T) {
	pLTable := &PLTable{
		Leader: Leader{
			LeaderId: LeaderId{
				Id: "1",
			},
		},
		PeerTable: map[string]Peer{
			"1": Peer{
				PeerId: PeerId{
					Id: "1",
				},
			},
		},
	}

	leader := &Leader{
		LeaderId: LeaderId{
			Id: "1",
		},
	}

	peerTable := map[string]Peer{
		"1": Peer{
			PeerId: PeerId{
				Id: "1",
			},
		},
	}

	p := NewPLTable(*leader, peerTable)

	assert.Equal(t, p, pLTable)
}

func TestPLTable_GetLeader(t *testing.T) {
	pLTable := &PLTable{
		Leader: Leader{
			LeaderId: LeaderId{
				Id: "1",
			},
		},
		PeerTable: map[string]Peer{
			"1": Peer{
				PeerId: PeerId{
					Id: "1",
				},
			},
		},
	}

	leader := &Leader{
		LeaderId: LeaderId{
			Id: "1",
		},
	}

	got, _ := pLTable.GetLeader()

	assert.Equal(t, *leader, got)
}

func TestPLTable_GetPeerTable(t *testing.T) {
	pLTable := &PLTable{
		Leader: Leader{
			LeaderId: LeaderId{
				Id: "1",
			},
		},
		PeerTable: map[string]Peer{
			"1": Peer{
				PeerId: PeerId{
					Id: "1",
				},
			},
		},
	}

	peerTable := map[string]Peer{
		"1": Peer{
			PeerId: PeerId{
				Id: "1",
			},
		},
	}

	pt, _ := pLTable.GetPeerTable()

	assert.Equal(t, pt, peerTable)
}

func TestPLTableService_GetPLTableFromCommand(t *testing.T) {
	pt := &PLTable{
		Leader: Leader{
			LeaderId: LeaderId{
				Id: "1",
			},
		},
		PeerTable: map[string]Peer{
			"1": Peer{
				PeerId: PeerId{
					Id: "1",
				},
			},
		},
	}

	byte, _ := common.Serialize(pt)

	c := &command.ReceiveGrpc{
		Body: byte,
	}

	pLTableService := &PLTableService{}
	extracted, _ := pLTableService.GetPLTableFromCommand(*c)

	assert.Equal(t, extracted, *pt)
}

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

package api_test

import (
	"testing"

	"github.com/it-chain/engine/common/test/mock"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/api"
	"github.com/it-chain/engine/p2p/infra/mem"
	"github.com/magiconair/properties/assert"
)

func TestLeaderApi_UpdateLeaderWithAddress(t *testing.T) {

	tests := map[string]struct {
		input struct {
			ipAddress string
		}
		output struct {
			leader p2p.Leader
		}
		err error
	}{
		"success": {
			input: struct{ ipAddress string }{ipAddress: "2"},
			output: struct{ leader p2p.Leader }{leader: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "2",
				},
			}},
			err: nil,
		},
		"no matching ipAddress test": {
			input: struct{ ipAddress string }{ipAddress: "234"},
			output: struct{ leader p2p.Leader }{
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
			},
			err: api.ErrNoMatchingPeerWithIpAddress,
		},
	}

	leaderApi := SetupLeaderApi(p2p.PLTable{
		Leader: p2p.Leader{
			LeaderId: p2p.LeaderId{
				Id: "1",
			},
		},
		PeerTable: map[string]p2p.Peer{
			"1": {
				PeerId: p2p.PeerId{
					Id: "1",
				},
				IpAddress: "1",
			},
			"2": {
				PeerId: p2p.PeerId{
					Id: "2",
				},
				IpAddress: "2",
			},
			"3": {
				PeerId: p2p.PeerId{
					Id: "3",
				},
				IpAddress: "3",
			},
		},
	})

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		err := leaderApi.UpdateLeaderWithAddress(test.input.ipAddress)

		//assert.Equal(t, mock.MakeFakePLTable().Leader, test.output.leader)

		assert.Equal(t, err, test.err)
	}
}

func TestLeaderApi_UpdateLeaderWithLargePeerTable(t *testing.T) {

	tests := map[string]struct {
		input struct {
			myPLTable       p2p.PLTable
			oppositePLTable p2p.PLTable
		}
		output struct {
			leader p2p.Leader
		}
	}{
		"success": {
			input: struct {
				myPLTable       p2p.PLTable
				oppositePLTable p2p.PLTable
			}{myPLTable: p2p.PLTable{
				Leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
				PeerTable: map[string]p2p.Peer{
					"1": {
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
				},
			}, oppositePLTable: p2p.PLTable{
				Leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "2",
					},
				},
				PeerTable: map[string]p2p.Peer{
					"1": {
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
					"2": {
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
					"3": {
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
					"4": {
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
				},
			}},
			output: struct{ leader p2p.Leader }{leader: p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: "2",
				},
			}},
		},
		"not updated with longer peer list case": {
			input: struct {
				myPLTable       p2p.PLTable
				oppositePLTable p2p.PLTable
			}{myPLTable: p2p.PLTable{
				Leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
				PeerTable: map[string]p2p.Peer{
					"1": {
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
					"2": {
						PeerId: p2p.PeerId{
							Id: "2",
						},
					},
					"3": {
						PeerId: p2p.PeerId{
							Id: "3",
						},
					},
				},
			}, oppositePLTable: p2p.PLTable{
				Leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "2",
					},
				},
				PeerTable: map[string]p2p.Peer{
					"1": {
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
					"2": {
						PeerId: p2p.PeerId{
							Id: "1",
						},
					},
				},
			}},
			output: struct {
				leader p2p.Leader
			}{
				leader: p2p.Leader{
					LeaderId: p2p.LeaderId{
						Id: "1",
					},
				},
			},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)

		leaderApi := SetupLeaderApi(test.input.myPLTable)

		leaderApi.UpdateLeaderWithLargePeerTable(test.input.oppositePLTable)

		leader, _ := leaderApi.PeerRepository.GetLeader()

		assert.Equal(t, leader, test.output.leader)
	}

}

func SetupLeaderApi(myPLTable p2p.PLTable) api.LeaderApi {

	peerRepository := mem.NewPeerReopository()

	peerRepository.SetLeader(myPLTable.Leader)

	for _, v := range myPLTable.PeerTable {

		peerRepository.Save(v)
	}

	eventService := mock.EventService{}

	eventService.PublishFunc = func(topic string, event interface{}) error {
		return nil
	}

	leaderApi := api.NewLeaderApi(&peerRepository, &eventService)

	return leaderApi
}

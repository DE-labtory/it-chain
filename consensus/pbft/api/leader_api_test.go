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

//func TestNewLeaderApi(t *testing.T) {
//	parliament := &pbft.Parliament{
//		Leader: &pbft.Leader{
//			LeaderId: "1",
//		},
//	}
//
//	eventService := &mock.EventService{}
//
//	leaderApi := api.NewLeaderApi(parliament, eventService)
//
//	assert.Equal(t, leaderApi.GetParliament(), parliament)
//}
//
//func TestLeaderApi_UpdateLeaderWithAddress(t *testing.T) {
//
//	leaderApi := SetLeaderApi()
//
//	// when
//	leaderApi.UpdateLeaderWithAddress("2")
//
//	//then
//	assert.Equal(t, leaderApi.GetLeader().LeaderId, "2")
//
//	//	when
//	err2 := leaderApi.UpdateLeaderWithAddress("4")
//
//	//	then
//	assert.Equal(t, err2, api.ErrNoMatchingPeerWithIpAddress)
//
//}
//
//func TestLeaderApi_GetLeader(t *testing.T) {
//	leaderApi := SetLeaderApi()
//
//	assert.Equal(t, leaderApi.GetLeader(), &pbft.Leader{LeaderId: "1"})
//}
//
//func TestLeaderApi_GetParliament(t *testing.T) {
//	leaderApi := SetLeaderApi()
//
//	parliament := &pbft.Parliament{
//		Leader: &pbft.Leader{
//			LeaderId: "1",
//		},
//		RepresentativeTable: map[string]*pbft.Representative{
//			"1": {
//				ID:        "1",
//				IpAddress: "1",
//			},
//			"2": {
//				ID:        "2",
//				IpAddress: "2",
//			},
//			"3": {
//				ID:        "3",
//				IpAddress: "3",
//			},
//		},
//	}
//
//	assert.Equal(t, leaderApi.GetParliament(), parliament)
//}
//
//func SetLeaderApi() *api.LeaderApi {
//	parliament := &pbft.Parliament{
//		Leader: &pbft.Leader{
//			LeaderId: "1",
//		},
//		RepresentativeTable: map[string]*pbft.Representative{
//			"1": {
//				ID:        "1",
//				IpAddress: "1",
//			},
//			"2": {
//				ID:        "2",
//				IpAddress: "2",
//			},
//			"3": {
//				ID:        "3",
//				IpAddress: "3",
//			},
//		},
//	}
//
//	eventService := &mock.EventService{}
//
//	eventService.PublishFunc = func(topic string, event interface{}) error {
//		return nil
//	}
//
//	leaderApi := api.NewLeaderApi(parliament, eventService)
//	return leaderApi
//}

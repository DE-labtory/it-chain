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
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/iLogger"
)

type LeaderEventHandler struct {
	electionApi *api.ElectionApi
}

func NewLeaderEventHandler(electionApi *api.ElectionApi) *LeaderEventHandler {
	return &LeaderEventHandler{
		electionApi: electionApi,
	}
}

func (l *LeaderEventHandler) HandlerLeaderDeletedEvent(_ event.LeaderDeleted) {
	iLogger.Infof(nil, "[PBFT] Leader Deleted, Start Elect Leader With RAFT")
	go l.electionApi.ElectLeaderWithRaft()
}

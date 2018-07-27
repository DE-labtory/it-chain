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

package api

import (
	"github.com/it-chain/engine/consensus/infra/adapter"
	"github.com/it-chain/engine/consensus"
)

type ConsensusApi struct {
	parliamentService adapter.ParliamentService
	consensusService consensus.Consensus
}

// todo : Event Sourcing 첨가

func (cApi ConsensusApi) StartConsensus(userId consensus.MemberId, block consensus.ProposedBlock) error {

	// 합의 시작!! 리더에 의해 시작 만약 블록이 생성되면 Consensus가 필요한지 따져야함
	// consensus를 시작한 멤버 아이디와, 제안된 블록으로 consensus를 만든다.
	// 합의 필요

	peerList, _ := cApi.parliamentService.RequestPeerList()
	if  cApi.parliamentService.IsNeedConsensus(){
		createdConsensus, err := consensus.CreateConsensus(peerList, block)

		if err != nil{
			print("error 발생")
			return nil
		}
		createdConsensus.Start()
		//TODO 다른 피어에게 메시지 전송

	}
	return nil
}

func (cApi ConsensusApi) ReceivePrePrepareMsg(msg consensus.PrePrepareMsg) error {
	// 검증하는 함수 if -> 검증.false == 수용 x
	// msg가 leader에게 온 것인지 검증
	// TODO message Service에 옮김 추가 검증 필요?
	lid, _ := cApi.parliamentService.RequestLeader()
	if lid.ToString() == msg.SenderId{
		// 검증 후 consensus Construct
		createdConsensus, err := consensus.ConstructConsensus(msg)
		if err != nil{
			print("construct consensus Err")
		}
		prepareMsg := consensus.NewPrepareMsg(createdConsensus)
		// TODO 모든피어에게 prepareMsg 보내야 함

		return nil
	}
	return nil
}

func (cApi ConsensusApi) ReceivePrepareMsg(msg consensus.PrepareMsg) error{
	return nil
}

func (cApi ConsensusApi) ReceiveCommitMsg(msg consensus.CommitMsg) error{
	return nil
}

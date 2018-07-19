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

package consensus

import "github.com/it-chain/midgard"

type CreateBlockCommand struct {
	midgard.CommandModel
	Block struct {
		Seal []byte
		Body []byte
	}
}

type SendPrePrepareMsgCommand struct {
	midgard.CommandModel
	PrePrepareMsg struct {
		ConsensusId    ConsensusId
		SenderId       string
		Representative []*Representative
		ProposedBlock  ProposedBlock
	}
}

type SendPrepareMsgCommand struct {
	midgard.CommandModel
	PrepareMsg struct {
		ConsensusId ConsensusId
		SenderId    string
		BlockHash   []byte
	}
}

type SendCommitMsgCommand struct {
	midgard.CommandModel
	CommitMsg struct {
		ConsensusId ConsensusId
		SenderId    string
	}
}

type SendGrpcMsgCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

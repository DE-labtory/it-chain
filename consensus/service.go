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

import "errors"

type PropagateService interface {
	BroadcastPrePrepareMsg(msg PrePrepareMsg) error
	BroadcastPrepareMsg(msg PrepareMsg) error
	BroadcastCommitMsg(msg CommitMsg) error
}

type ConfirmService interface {
	ConfirmBlock(block ProposedBlock) error
}

// 연결된 peer 중에서 consensus 에 참여할 representative 들을 선출
func Elect(parliament []MemberId) ([]*Representative, error) {
	representatives := make([]*Representative, 0)

	if len(parliament) == 0 {
		return []*Representative{}, errors.New("No parliament member.")
	}

	for _, peerId := range parliament {
		representatives = append(representatives, NewRepresentative(peerId.ToString()))
	}

	return representatives, nil
}

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

import "github.com/it-chain/midgard"

//다른 Peer에게 Message수신 command
type GrpcReceiveCommand struct {
	midgard.CommandModel
	Body         []byte
	ConnectionID string
	Protocol     string
}

//다른 Peer에게 Message전송 command
type GrpcDeliverCommand struct {
	midgard.CommandModel
	Recipients []string //connectionId
	Body       []byte
	Protocol   string
}

//Connection 생성 command
type ConnectionCreateCommand struct {
	midgard.CommandModel
	Address string
}

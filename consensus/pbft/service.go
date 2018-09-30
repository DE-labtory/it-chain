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

package pbft

import (
	"errors"
	"math/rand"
	"time"
)

var ErrNoParliamentMember = errors.New("No parliament member.")

// 연결된 peer 중에서 consensus 에 참여할 representative 들을 선출
func Elect(parliament []Representative) ([]Representative, error) {
	if len(parliament) == 0 {
		return []Representative{}, ErrNoParliamentMember
	}
	return parliament, nil
}

func GenRandomInRange(min, max int) int {

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

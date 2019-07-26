/*
 * Copyright 2018 DE-labtory
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

package batch

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//todo return error 했을 경우 test할 방법이 없음
func TestTimeoutBatcher_Run(t *testing.T) {

	wg := sync.WaitGroup{}

	counter := 0
	//given
	tests := map[string]struct {
		input struct {
			taskFunc TaskFunc
			duration time.Duration
		}
		err error
	}{
		"success": {
			input: struct {
				taskFunc TaskFunc
				duration time.Duration
			}{
				taskFunc: func() error {
					if counter == 0 {
						wg.Done()
						fmt.Println("success done")
					}

					counter++

					return nil
				},
				duration: time.Second * 3,
			},
			err: nil,
		},
	}

	batcher := GetTimeOutBatcherInstance()
	wg.Add(1)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		counter = 0
		quit := batcher.Run(test.input.taskFunc, test.input.duration)
		wg.Wait()
		quit <- struct{}{}
	}
}

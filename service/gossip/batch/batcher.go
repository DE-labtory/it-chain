/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package batch

import "it-chain/common"

var logger = common.GetLogger("batch")

type Batcher interface {

	// Add adds a message to be batched
	Add(interface{})

	// Stop stops the component
	Stop()

	// Size returns the amount of pending messages to be emitted
	Size() int
}




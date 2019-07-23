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

package mock_test

import (
	"testing"
	"time"

	"github.com/DE-labtory/it-chain/common/command"
	"github.com/DE-labtory/it-chain/common/mock"
)

func TestEventService_Publish(t *testing.T) {
	eventService := mock.NewEventService("1", func(processId string, topic string, event interface{}) error {
		return nil
	})

	eventService.SetDelayTime(5 * time.Millisecond)

	event := command.DeliverGrpc{
		MessageId: "1",
	}

	eventService.Publish("message.deliver", event)
}

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

package mock

import "github.com/DE-labtory/midgard"

type EventRepository struct {
	LoadFunc  func(aggregate midgard.Aggregate, aggregateID string) error
	SaveFunc  func(aggregateID string, events ...midgard.Event) error
	CloseFunc func()
}

func (er EventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return er.LoadFunc(aggregate, aggregateID)
}
func (er EventRepository) Save(aggregateID string, events ...midgard.Event) error {
	return er.SaveFunc(aggregateID, events...)
}
func (er EventRepository) Close() {
	er.CloseFunc()
}

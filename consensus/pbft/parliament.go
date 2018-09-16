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

type Parliament struct {
	Leader              *Leader
	RepresentativeTable map[string]*Representative
}

func NewParliament() *Parliament {
	return &Parliament{
		Leader:              &Leader{},
		RepresentativeTable: make(map[string]*Representative),
	}
}

type Leader struct {
	LeaderId string
}

func (l Leader) GetID() string {
	return l.LeaderId
}

type Representative struct {
	ID        string
	IpAddress string
}

func (r Representative) GetID() string {
	return string(r.ID)
}

func NewRepresentative(ID string) *Representative {
	return &Representative{ID: ID}
}

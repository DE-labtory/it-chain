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

package model

// it-chain의 공통적이고 중요한 설정을 담는 구조체이다.
type EngineConfiguration struct {
	KeyPath              string
	Mode                 string
	Amqp                 string
	BootstrapNodeAddress string
}

func NewEngineConfiguration() EngineConfiguration {
	return EngineConfiguration{
		KeyPath:              ".it-chain/",
		Mode:                 "solo",
		BootstrapNodeAddress: "127.0.0.1:5555",
		Amqp:                 "amqp://guest:guest@localhost:5672/",
	}
}

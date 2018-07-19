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

package conf

import (
	"os"
	"sync"

	"github.com/it-chain/engine/conf/model"
	"github.com/spf13/viper" //viper는 go 어플리케이션의 각종 설정을 담당하는 lib이다.
	// 각종 형태의 설정파일을 찾고, 로드하는 것이 주 역할이다.
)

// it-chain 설정을 저장하는 구조체에 대한 포인터 instance를 선언한다.
var instance = &Configuration{}

// 특정 함수를 한번만 수행하기 위한 sync.Once 를 once로 선언하고 필요한 경우
// once.Do(func(){}) 로 호출하여 사용하도록 한다.
var once sync.Once

// it-chain 에 필요한 각종 설정을 저장하는 구조체이다.
type Configuration struct {
	configName  string
	Engine      model.EngineConfiguration
	Txpool      model.TxpoolConfiguration
	Consensus   model.ConsensusConfiguration
	Blockchain  model.BlockChainConfiguration
	Peer        model.PeerConfiguration
	Icode       model.ICodeConfiguration
	GrpcGateway model.GrpcGatewayConfiguration
	ApiGateway  model.ApiGatewayConfiguration
}

// it-chain의 각종 설정을 받아온다.
func SetConfigName(name string) {
	instance.configName = name
}
func GetConfiguration() *Configuration {
	if instance.configName == "" {
		instance.configName = "config"
	}
	// 최초로 go application의 configuration을 해당 파일을 통해 설정한다.
	once.Do(func() {

		// instance를 it-chain 설정에 관한 구조체의 포인터로 지정한다.

		// Go language의 환경변수와 내부 디렉터리 구조를 통해 config 파일이 저장된 위치와 파일명을 잡아준다.
		path := os.Getenv("GOPATH") + "/src/github.com/it-chain/engine/conf"
		viper.SetConfigName(instance.configName)
		viper.AddConfigPath(path)

		if err := viper.ReadInConfig(); err != nil {
			panic("cannot read config")
		}
		err := viper.Unmarshal(&instance)
		if err != nil {
			panic("error in read config")
		}
	})

	// it-chain의 설정내용을 반환한다.
	return instance
}

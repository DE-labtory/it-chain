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

package conf

import (
	"os"
	"sync"

	"github.com/DE-labtory/it-chain/conf/model"
	"github.com/spf13/viper" //viper는 go 어플리케이션의 각종 설정을 담당하는 lib이다.
)

// it-chain 설정을 저장하는 구조체에 대한 포인터 instance를 선언한다.
var instance = &Configuration{}

// 특정 함수를 한번만 수행하기 위한 sync.Once 를 once로 선언하고 필요한 경우
// once.Do(func(){}) 로 호출하여 사용하도록 한다.
var once sync.Once

var confPath = os.Getenv("GOPATH") + "/src/github.com/DE-labtory/it-chain/conf/config.yaml"

// it-chain 에 필요한 각종 설정을 저장하는 구조체이다.
type Configuration struct {
	Engine      model.EngineConfiguration
	Txpool      model.TxpoolConfiguration
	Consensus   model.ConsensusConfiguration
	Blockchain  model.BlockChainConfiguration
	Peer        model.PeerConfiguration
	Icode       model.ICodeConfiguration
	GrpcGateway model.GrpcGatewayConfiguration
	ApiGateway  model.ApiGatewayConfiguration
	Docker      model.DockerConfiguration
}

// EngineConfiguration Mode용
var modeConsts = [...]string{
	"solo",
	"test",
	"pbft",
}

// GOPATH 설정유무 확인, conf package 호출 시 최초 실행되는 func
func init() {
	if os.Getenv("GOPATH") == "" {
		panic("Need to set GOPATH")
	}
}

// it-chain의 conf path를 받아온다.
func SetConfigPath(path string) {
	confPath = path
}

func GetConfiguration() *Configuration {

	once.Do(func() {
		viper.SetConfigFile(confPath)
		if err := viper.ReadInConfig(); err != nil {
			panic("cannot read config")
		}
		err := viper.Unmarshal(&instance)
		if err != nil {
			panic("error in read config")
		}

		if !HasValidMode(instance) {
			panic("NewEngineConfiguration mode is wrong.")
		}
	})
	return instance
}

func HasValidMode(c *Configuration) bool {

	for _, modeConst := range modeConsts {
		if instance.Engine.Mode == modeConst {
			return true
		}
	}
	return false

}

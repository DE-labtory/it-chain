package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"flag"

	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/conf/model"
	"github.com/it-chain/it-chain-Engine/conf/model/common"
	"gopkg.in/yaml.v2"
)

type test struct {
	A   string
	B   string
	Inn inner
}
type inner struct {
	Innerint int
}

func main() {
	path, _ := os.Getwd()
	configName := flag.String("name", "config", "config file name")
	flag.Parse()
	if _, err := os.Stat(path + "/" + *configName + ".yaml"); err == nil {
		for i := 0; ; i++ {
			if _, err := os.Stat(path + "/" + *configName + "_bak" + strconv.Itoa(i) + ".yaml"); os.IsNotExist(err) {
				os.Rename(path+"/"+*configName+".yaml", path+"/"+*configName+"_bak"+strconv.Itoa(i)+".yaml")
				break
			}
		}
	}

	confInfo := conf.Configuration{
		Common:         common.NewCommonConfiguration(),
		Txpool:         model.NewTxpoolConfiguration(),
		Consensus:      model.NewConsensusConfiguration(),
		Blockchain:     model.NewBlockChainConfiguration(),
		Peer:           model.NewPeerConfiguration(),
		Authentication: model.NewAuthenticationConfiguration(),
		Icode:          model.NewIcodeConfiguration(),
		GrpcGateway:    model.NewGrpcGatewayConfiguration(),
	}

	output, _ := yaml.Marshal(&confInfo)
	err := ioutil.WriteFile(path+"/"+*configName+".yaml", output, 0644)

	if err != nil {
		fmt.Println(err.Error())
		panic("Error in generate config file")
	}
	println("success to generate config file")
}

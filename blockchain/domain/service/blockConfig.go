package service

import (
	"go/build"
	"io/ioutil"
	"os"

	_ "github.com/it-chain/it-chain-Engine/common"
	"github.com/spf13/viper"
)

func ConfigFromJson(filename string) ([]uint8, error) {

	enginePath := build.Default.GOPATH + "/src/github.com/it-chain/it-chain-Engine/"
	folderPath := viper.GetString("genesisblock.defaultPath") + "/"
	filePath := enginePath + folderPath + filename
	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}

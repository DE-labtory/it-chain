package service

import (
	"io/ioutil"
	"os"

	_ "github.com/it-chain/it-chain-Engine/common"
)

func ConfigFromJson(filePath string) ([]uint8, error) {

	//enginePath := build.Default.GOPATH + "/src/github.com/it-chain/it-chain-Engine/"
	//folderPath := viper.GetString("genesisconf.defaultPath") + "/"
	//filePath := enginePath + folderPath + filename
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

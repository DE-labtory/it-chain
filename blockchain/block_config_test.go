package blockchain

import (
	"go/build"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFromJson(t *testing.T) {
	genesisconfPath := build.Default.GOPATH + "/src/github.com/it-chain/it-chain-Engine/.it-chain/genesisconf/"
	rightFilePath := genesisconfPath + "GenesisBlockConfig.json"
	wrongFilePath := genesisconfPath + "WrongFileName.json"
	_, err1 := ConfigFromJson(rightFilePath)
	assert.NoError(t, err1)
	_, err2 := ConfigFromJson(wrongFilePath)
	assert.Error(t, err2)
}

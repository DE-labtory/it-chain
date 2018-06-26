package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/it-chain/yggdrasill/impl"
	"github.com/stretchr/testify/assert"
)

func TestConfigFromJson(t *testing.T) {

	genesisFilePath := "./GenesisBlockConfig.json"
	wrongFilePath := "./WrongFileName.json"
	var tempJson impl.DefaultBlock

	err := json.Unmarshal(GenesisBlockConfigJson, &tempJson)
	assert.NoError(t, err)

	GenesisBlockConfigByte, _ := json.Marshal(tempJson)
	err = ioutil.WriteFile(genesisFilePath, GenesisBlockConfigByte, 0644)
	assert.NoError(t, err)

	defer os.Remove(genesisFilePath)

	_, err1 := ConfigFromJson(genesisFilePath)
	assert.NoError(t, err1)
	_, err2 := ConfigFromJson(wrongFilePath)
	assert.Error(t, err2)
}

package blockchain

import (
	"testing"

	"os"

	"io/ioutil"

	"encoding/json"

	"github.com/it-chain/yggdrasill/impl"
	"github.com/stretchr/testify/assert"
)

func TestConfigFromJson(t *testing.T) {
	genesisFilePath := "./GenesisBlockConfig.json"
	defer os.Remove(genesisFilePath)
	wrongFilePath := "./WrongFileName.json"
	GenesisBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)
	var tempJson impl.DefaultBlock
	_ = json.Unmarshal(GenesisBlockConfigJson, &tempJson)
	GenesisBlockConfigByte, _ := json.Marshal(tempJson)
	_ = ioutil.WriteFile(genesisFilePath, GenesisBlockConfigByte, 0644)
	_, err1 := ConfigFromJson(genesisFilePath)
	assert.NoError(t, err1)
	_, err2 := ConfigFromJson(wrongFilePath)
	assert.Error(t, err2)
}

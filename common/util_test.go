package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"it-chain/service/blockchain"
)

func TestCreateDirIfMissing(t *testing.T){

	dirPath := "./test_path"

	CreateDirIfMissing(dirPath)
	assert.DirExists(t,dirPath)

	//clean up
	os.Remove(dirPath)
}

func TestDirEmpty(t *testing.T) {

	dirPath := "./test_path/"

	os.MkdirAll(path.Dir(dirPath), 0755)

	isExist, err := DirEmpty(dirPath)

	if err != nil{
		//error
	}

	assert.True(t,isExist)

	//clean up
	os.Remove(dirPath)
}

func TestSerialize(t *testing.T) {
	block := &blockchain.Block{}

	serialized, err := Serialize(block)
	assert.NoError(t, err)

	deserialized, err := Deserialize(serialized, &blockchain.Block{})
	assert.NoError(t, err)
	assert.Equal(t, block, deserialized)
}
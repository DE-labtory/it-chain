package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"time"
)

type TestStrunct struct {
	MemberString string
	MemberInt    int
	MemberTime   time.Time
	MemberByte   []byte
}

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
	testStruct := &TestStrunct{}

	serialized, err := Serialize(testStruct)
	assert.NoError(t, err)

	deserialized, err := Deserialize(serialized, &TestStrunct{})
	assert.NoError(t, err)
	assert.Equal(t, testStruct, deserialized)
}
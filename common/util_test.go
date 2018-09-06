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

package common

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	MemberString string
	MemberInt    int
	MemberTime   time.Time
	MemberByte   []byte
}

func TestCreateDirIfMissing(t *testing.T) {

	dirPath := "./test_path"

	CreateDirIfMissing(dirPath)
	assert.DirExists(t, dirPath)

	//clean up
	os.Remove(dirPath)
}

func TestDirEmpty(t *testing.T) {

	dirPath := "./test_path/"

	os.MkdirAll(path.Dir(dirPath), 0755)

	isExist, err := DirEmpty(dirPath)

	if err != nil {
		//error
	}

	assert.True(t, isExist)

	//clean up
	os.Remove(dirPath)
}

func TestSerialize(t *testing.T) {
	testStruct := &TestStruct{MemberString: "test"}

	serialized, err := Serialize(testStruct)
	assert.NoError(t, err)

	data := &TestStruct{}
	err = Deserialize(serialized, data)
	assert.NoError(t, err)
	assert.Equal(t, testStruct, data)
}

func TestRelativeToAbsolutePath(t *testing.T) {

	testfile1 := "./util.go"
	testabsresult1, err := filepath.Abs(testfile1)
	assert.NoError(t, err)
	testabs1, err := RelativeToAbsolutePath(testfile1)

	assert.NoError(t, err)
	assert.Equal(t, testabs1, testabsresult1)

	testfile2 := "../README.md"
	testabsresult2, err := filepath.Abs(testfile2)
	assert.NoError(t, err)

	testabs2, err := RelativeToAbsolutePath(testfile2)

	assert.NoError(t, err)
	assert.Equal(t, testabs2, testabsresult2)

	// 남의 홈패스에 뭐가있는지 알길이 없으니 하나 만들었다 지움
	usr, err := user.Current()
	assert.NoError(t, err)

	testfile3 := usr.HomeDir + "/test.txt"

	_, err = os.Stat(usr.HomeDir)
	if os.IsNotExist(err) {
		file, err := os.Create(testfile3)
		assert.NoError(t, err)
		defer file.Close()
	}

	err = ioutil.WriteFile(testfile3, []byte("test"), os.ModePerm)
	assert.NoError(t, err)

	testfile4 := "~/test.txt"

	testabs3, err := RelativeToAbsolutePath(testfile4)
	assert.NoError(t, err)
	assert.Equal(t, testfile3, testabs3)

	err = os.Remove(testfile3)
	assert.NoError(t, err)
}

// TODO
func TestRelativeToAbsolutePath_WhenGivenPathIsAbsolute(t *testing.T) {

}

// TODO
func TestRelativeToAbsolutePath_WhenGivenPathWithOnlyName(t *testing.T) {

}

// TODO
func TestRelativeToAbsolutePath_WhenGIvenPathIsEmpty(t *testing.T) {

}

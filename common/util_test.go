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

package common_test

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/it-chain/engine/common"

	"regexp"

	"fmt"

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

	common.CreateDirIfMissing(dirPath)
	assert.DirExists(t, dirPath)

	//clean up
	os.Remove(dirPath)
}

func TestDirEmpty(t *testing.T) {

	dirPath := "./test_path/"

	os.MkdirAll(path.Dir(dirPath), 0755)

	isExist, err := common.DirEmpty(dirPath)

	if err != nil {
		//error
	}

	assert.True(t, isExist)

	//clean up
	os.Remove(dirPath)
}

func TestSerialize(t *testing.T) {
	testStruct := &TestStruct{MemberString: "test"}

	serialized, err := common.Serialize(testStruct)
	assert.NoError(t, err)

	data := &TestStruct{}
	err = common.Deserialize(serialized, data)
	assert.NoError(t, err)
	assert.Equal(t, testStruct, data)
}

func TestRelativeToAbsolutePath(t *testing.T) {

	testfile1 := "./util.go"
	testabsresult1, err := filepath.Abs(testfile1)
	assert.NoError(t, err)
	testabs1, err := common.RelativeToAbsolutePath(testfile1)

	assert.NoError(t, err)
	assert.Equal(t, testabs1, testabsresult1)

	testfile2 := "../README.md"
	testabsresult2, err := filepath.Abs(testfile2)
	assert.NoError(t, err)

	testabs2, err := common.RelativeToAbsolutePath(testfile2)

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

	testabs3, err := common.RelativeToAbsolutePath(testfile4)
	assert.NoError(t, err)
	assert.Equal(t, testfile3, testabs3)

	err = os.Remove(testfile3)
	assert.NoError(t, err)
}

func TestRelativeToAbsolutePath_WhenGivenPathIsAbsolute(t *testing.T) {
	sshPath := "/iAmRoot"

	absPath, err := common.RelativeToAbsolutePath(sshPath)

	assert.NoError(t, err)
	assert.Equal(t, sshPath, absPath)
}

func TestRelativeToAbsolutePath_WhenGivenPathWithOnlyName(t *testing.T) {
	sshPath := "test-dir"

	absPath, err := common.RelativeToAbsolutePath(sshPath)
	currentPath, _ := filepath.Abs(".")

	assert.NoError(t, err)
	assert.Equal(t, path.Join(currentPath, sshPath), absPath)
}

func TestRelativeToAbsolutePath_WhenGivenPathIsEmpty(t *testing.T) {
	sshPath := ""

	absPath, err := common.RelativeToAbsolutePath(sshPath)

	assert.Equal(t, nil, err)
	assert.Equal(t, "", absPath)
}

func TestCheckRegex(t *testing.T) {
	r, err := regexp.Compile("block")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r.MatchString("block.created"))
}

func TestFindEarliestString(t *testing.T) {
	larger := common.FindEarliestString([]string{"a", "b", "c"})

	assert.Equal(t, "a", larger)
}

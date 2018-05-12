/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path"
	"sort"
	"strings"
)

var logger = GetLogger("util.go")

func CreateDirIfMissing(dirPath string) error {

	if !strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath + "/"
	}

	//logger.Debugf("CreateDirIfMissing [%s]", dirPath)

	err := os.MkdirAll(path.Dir(dirPath), 0755)
	if err != nil {
		logger.Debugf("Error while creating dir [%s]", dirPath)
		return err
	}

	return nil
}

// DirEmpty returns true if the dir at dirPath is empty
func DirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		logger.Debugf("Error while opening dir [%s]: %s", dirPath, err)
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func ComputeSHA256(data []string) string {
	sort.Strings(data)
	arg := strings.Join(data, ",")
	hash := sha256.New()
	hash.Write([]byte(arg))
	return hex.EncodeToString(hash.Sum(nil))

}

/**
gob encoder로 인코딩했을 때 문제점
1. empty slice(make 로 생성한거) 가 디코딩하면 nil 로 디코딩 됨.
ㄴ json marshal로 바꾸면서 해결
2. time.Time 값들은 뒤에 monotonic 파트가 없어짐.
2번은 문제가 안 될수도 있는데 테스트 실패의 원인..
*/
func Serialize(object interface{}) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		panic(fmt.Sprintf("Error encoding : %s", err))
	}
	return data, nil
}

func Deserialize(serializedBytes []byte, object interface{}) error {
	if len(serializedBytes) == 0 {
		return nil
	}
	err := json.Unmarshal(serializedBytes, object)
	if err != nil {
		panic(fmt.Sprintf("Error decoding : %s", err))
	}
	return err
}

func CryptoRandomGeneration(min int64, max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max+1-min))
	ret := n.Int64() + min
	return ret
}

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
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/it-chain/engine/common/command"
	"github.com/rs/xid"
)

func CreateDirIfMissing(dirPath string) error {

	if !strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath + "/"
	}

	//logger.Debugf("CreateDirIfMissing [%s]", dirPath)

	err := os.MkdirAll(path.Dir(dirPath), 0755)
	if err != nil {
		return err
	}

	return nil
}

// DirEmpty returns true if the dir at dirPath is empty
func DirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
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

// encoding/json 패키지에 대한 설명
// json data를 json으로 인코딩 하기 위한 encoding/json 패키지의 Marshal() 을 사용한다.
// json.Marshal(_struct_) 을 하면 (json 인코딩 바이트배열, 에러객체) 를 리턴한다.
// decode를 위해서는 json.Unmarshal(_jsonBytes_, 받을 구조체 포인터) 를 사용한다.

func Serialize(object interface{}) ([]byte, error) { //모든 stuct 받기 위해 interface{} 타입의 입력으로 선언
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

// absolute path로 변경하기
// TODO: rename func
// TODO: Refactoring
func RelativeToAbsolutePath(rpath string) (string, error) {
	if rpath == "" {
		return rpath, nil
	}

	absolutePath := ""

	// 1. ./ ../ 경우
	if strings.Contains(rpath, "./") {
		abs, err := filepath.Abs(rpath)
		if err != nil {
			return rpath, err
		}
		return abs, nil
	}

	// 2. ~/ 홈폴더 경우
	if strings.Contains(rpath, "~") {
		i := strings.Index(rpath, "~") // 처음 나온 ~만 반환

		if i > -1 {
			pathRemain := rpath[i+1:]
			// user home 얻기
			usr, err := user.Current()
			if err != nil {
				return rpath, err
			}
			return path.Join(usr.HomeDir, pathRemain), nil

		} else {
			return rpath, nil
		}
	}

	if string(rpath[0]) == "/" {
		return rpath, nil
	}

	if string(rpath[0]) != "." && string(rpath[0]) != "/" {
		currentPath, err := filepath.Abs(".")
		if err != nil {
			return rpath, err
		}

		return path.Join(currentPath, rpath), nil
	}

	return absolutePath, nil

}

func FindEarliestString(list []string) string {
	largerOne := list[0]
	for _, v := range list {
		if strings.Compare(largerOne, v) > 0 {
			largerOne = v
		}
	}

	return largerOne
}

func CreateGrpcDeliverCommand(protocol string, body interface{}) (command.DeliverGrpc, error) {

	data, err := Serialize(body)

	if err != nil {
		return command.DeliverGrpc{}, err
	}

	return command.DeliverGrpc{
		MessageId:     xid.New().String(),
		RecipientList: make([]string, 0),
		Body:          data,
		Protocol:      protocol,
	}, err
}

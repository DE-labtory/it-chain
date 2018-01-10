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
	"strings"
	"os"
	"path"
	"io"
	"crypto/sha256"
)

var logger = GetLogger("util.go")

func CreateDirIfMissing(dirPath string) (error){

	if !strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath + "/"
	}

	logger.Debugf("CreateDirIfMissing [%s]", dirPath)

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

func ComputeSHA256(data []byte) (hash [32]uint8) {
	hash = sha256.Sum256(data)
	return hash
}
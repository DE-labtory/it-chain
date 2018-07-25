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

package logger_test

import (
	"testing"

	"path/filepath"

	"os"

	"github.com/it-chain/engine/common/logger"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	logger.SetToDebug()
	logger.EnableFileLogger(false, "")
	logger.Debug(nil, "debug test")
}

func TestError(t *testing.T) {
	logger.EnableFileLogger(false, "")
	logger.Errorf(&logger.Fields{"filedTest": "good"}, "%s good?", "testing is")
}

func TestEnableFileLogger(t *testing.T) {
	os.RemoveAll("./.tmp")
	absPath, _ := filepath.Abs("./.tmp/tmp.log")
	defer os.RemoveAll("./.tmp")
	err := logger.EnableFileLogger(true, absPath)
	assert.NoError(t, err)
	logger.Error(nil, "hahaha")
}

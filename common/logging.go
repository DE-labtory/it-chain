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
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Log = logrus.New()

//todo config을 통해 log 설정을 수행한다. (level 등등)
func init() {
	customFormatter := new(prefixed.TextFormatter)
	customFormatter.FullTimestamp = true

	Log.Formatter = customFormatter
	Log.Level = logrus.DebugLevel
}

func GetLogger(name string) *logrus.Entry {
	return Log.WithFields(logrus.Fields{
		"File": name,
	})
}

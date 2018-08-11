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
 *
 */

package logger

import (
	"os"

	"fmt"

	"runtime"

	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Fields = logrus.Fields

var (
	stdLogger  = logrus.New()
	fileLogger = logrus.New()
)

func init() {
	initStdLogger(stdLogger)
	fileLogger = nil
}

// default std logger is enabled
func EnableStd(enable bool) {
	if enable && stdLogger == nil {
		stdLogger = logrus.New()
		initStdLogger(stdLogger)
	} else if enable == false {
		stdLogger = nil
	}
}

//default file logger is disabled
func EnableFileLogger(enable bool, savePath string) error {
	if enable {
		fileLogger = logrus.New()
		err := initFileLogger(fileLogger, savePath)
		if err != nil {
			fileLogger = nil
		}
		return err
	} else {
		fileLogger = nil
		return nil
	}
}

func Debug(fields *Fields, message string) {
	if fields == nil {
		fields = &Fields{}
	}
	if stdLogger != nil {
		stdLogger.WithFields(*fields).Debug(message)
	}
	if fileLogger != nil {
		fileLogger.WithFields(*fields).Debug(message)
	}
}
func Info(fields *Fields, message string) {
	if fields == nil {
		fields = &Fields{}
	}
	if stdLogger != nil {
		stdLogger.WithFields(*fields).Info(message)
	}
	if fileLogger != nil {
		fileLogger.WithFields(*fields).Info(message)
	}
}
func Warn(fields *Fields, message string) {
	if fields == nil {
		fields = &Fields{}
	}
	if stdLogger != nil {
		fields = addLineInfo(fields)
		stdLogger.WithFields(*fields).Warn(message)
	}
	if fileLogger != nil {
		fields = addLineInfo(fields)
		fileLogger.WithFields(*fields).Warn(message)
	}
}
func Fatal(fields *Fields, message string) {
	if fields == nil {
		fields = &Fields{}
	}
	if stdLogger != nil {
		fields = addLineInfo(fields)
		stdLogger.WithFields(*fields).Fatal(message)
	}
	if fileLogger != nil {
		fields = addLineInfo(fields)
		fileLogger.WithFields(*fields).Fatal(message)
	}
}
func Error(fields *Fields, message string) {
	if fields == nil {
		fields = &Fields{}
	}
	if stdLogger != nil {
		fields = addLineInfo(fields)
		stdLogger.WithFields(*fields).Error(message)
	}
	if fileLogger != nil {
		fields = addLineInfo(fields)
		fileLogger.WithFields(*fields).Error(message)
	}
}

func Panic(fields *Fields, message string) {
	if fields == nil {
		fields = &Fields{}
	}
	if stdLogger != nil {
		fields = addLineInfo(fields)
		stdLogger.WithFields(*fields).Panic(message)
	}
	if fileLogger != nil {
		fields = addLineInfo(fields)
		fileLogger.WithFields(*fields).Panic(message)
	}
}

// formatting support
func Debugf(fields *Fields, format string, args ...interface{}) {
	Debug(fields, fmt.Sprintf(format, args...))
}

func Infof(fields *Fields, format string, args ...interface{}) {
	Info(fields, fmt.Sprintf(format, args...))
}
func Warnf(fields *Fields, format string, args ...interface{}) {
	Warn(fields, fmt.Sprintf(format, args...))
}
func Fatalf(fields *Fields, format string, args ...interface{}) {
	Fatal(fields, fmt.Sprintf(format, args...))
}
func Errorf(fields *Fields, format string, args ...interface{}) {
	Error(fields, fmt.Sprintf(format, args...))
}
func Panicf(fields *Fields, format string, args ...interface{}) {
	Panic(fields, fmt.Sprintf(format, args...))
}

func SetToDebug() {
	if stdLogger != nil {
		stdLogger.Level = logrus.DebugLevel
	}
	if fileLogger != nil {
		fileLogger.Level = logrus.DebugLevel
	}
}

func initStdLogger(logger *logrus.Logger) {
	logger.SetOutput(os.Stdout)
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:      true,
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
	}
}

func initFileLogger(logger *logrus.Logger, savePath string) error {

	if _, err := os.Stat(savePath); err != nil {
		err = os.MkdirAll(filepath.Dir(savePath), 0777)
		if err != nil {
			return err
		}
		_, err = os.Create(savePath)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(savePath, os.O_APPEND|os.O_WRONLY, 0777)
	if err == nil {
		logger.Formatter = &logrus.JSONFormatter{}
		logger.Out = file
	} else {
		return err
	}

	return nil
}

func addLineInfo(fields *Fields) *Fields {
	pc, _, _, _ := runtime.Caller(2)
	dataField := Fields{"cause": runtime.FuncForPC(pc).Name()}
	data := make(Fields, len(dataField)+len(*fields))
	for k, v := range dataField {
		data[k] = v
	}
	for k, v := range *fields {
		data[k] = v
	}
	return &data
}

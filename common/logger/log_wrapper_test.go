package logger

import (
	"testing"

	"path/filepath"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	EnableFileLogger(false, "")
	Debug(nil, "debug test")
}

func TestError(t *testing.T) {
	EnableFileLogger(false, "")
	Errorf(&Fields{"filedTest": "good"}, "%s good?", "testing is")
}

func TestEnableFileLogger(t *testing.T) {
	os.RemoveAll("./.tmp")
	absPath, _ := filepath.Abs("./.tmp/tmp.log")
	defer os.RemoveAll("./.tmp")
	err := EnableFileLogger(true, absPath)
	assert.NoError(t, err)
	Error(nil, "hahaha")

}

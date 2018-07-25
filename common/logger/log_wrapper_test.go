package logger

import "testing"

func TestDebug(t *testing.T) {
	EnableFileLogger(false, "")
	Debug(nil, "debug test")
}

func TestError(t *testing.T) {
	EnableFileLogger(false, "")
	Errorf(&Fields{"filedTest": "good"}, "%s good?", "testing is")
}

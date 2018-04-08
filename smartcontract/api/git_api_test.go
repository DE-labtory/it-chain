package api

import (
	"os"
	"testing"
)

func TestClone(t *testing.T) {

	var sshPath = "git@github.com:it-chain/tesseract.git"
	os.RemoveAll(".tmp")
	Clone(sshPath)
}

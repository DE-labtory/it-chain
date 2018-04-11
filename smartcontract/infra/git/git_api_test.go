package git

import (
	"os"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

var gitUrl = "git@github.com:it-chain/tesseract.git"

func TestClone(t *testing.T) {

	os.RemoveAll("./.tmp")
	api := NewGitApi()

	sc, err := api.Clone(gitUrl)

	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, "./.tmp/"+"tesseract.git", sc.Path)
	assert.Equal(t, gitUrl, sc.GitUrl)
	assert.Equal(t, "tesseract.git", sc.RepositoryName)
}

func TestGetNameFromGitUrl(t *testing.T) {
	fmt.Print(getNameFromGitUrl(gitUrl))
}

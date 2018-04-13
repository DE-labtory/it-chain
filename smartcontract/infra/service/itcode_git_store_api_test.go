package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var gitUrl = "git@github.com:it-chain/tesseract.git"

func TestClone(t *testing.T) {

	os.RemoveAll("./.tmp")
	defer os.RemoveAll("./.tmp")
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
	name := getNameFromGitUrl(gitUrl)

	assert.Equal(t, "tesseract.git", name)
}

func TestDirExist(t *testing.T) {

	p := "tmp"
	wp := "tmp2"
	defer os.RemoveAll(p)

	err := os.MkdirAll(p, 0755)

	assert.NoError(t, err)
	assert.True(t, dirExists(p))
	assert.False(t, dirExists(wp))
}

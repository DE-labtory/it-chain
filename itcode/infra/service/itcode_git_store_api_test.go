package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
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

	assert.Equal(t, "./.tmp/"+"tesseract", sc.Path)
	assert.Equal(t, gitUrl, sc.GitUrl)
	assert.Equal(t, "tesseract", sc.RepositoryName)
}

func TestChangeRemote(t *testing.T) {

	//given
	os.RemoveAll("./.tmp")
	defer os.RemoveAll("./.tmp")
	api := NewGitApi()
	itCode, err := api.Clone(gitUrl)
	assert.NoError(t, err)

	//when
	err = api.ChangeRemote(itCode.Path, "https://github.com/steve-buzzni"+"/"+itCode.RepositoryName)
	assert.NoError(t, err)

	//then
	r, err := git.PlainOpen(itCode.Path)
	assert.NoError(t, err)
	remote, err := r.Remote(git.DefaultRemoteName)
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/steve-buzzni"+"/"+itCode.RepositoryName, remote.Config().URLs[0])
}

func TestGitApi_Push(t *testing.T) {

	repoName := "test-chain"
	api := NewGitApi()
	itCode, err := api.Clone(gitUrl)
	assert.NoError(t, err)

	_, _, err := CreateRepository(repoName)
	assert.NoError(t, err)

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

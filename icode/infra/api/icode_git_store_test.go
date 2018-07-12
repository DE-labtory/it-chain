package api_test

import (
	"os"
	"testing"

	"errors"

	"strings"

	"github.com/google/go-github/github"
	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/it-chain/it-chain-Engine/icode/infra/api"
	"github.com/stretchr/testify/assert"
)

func TestICodeGitStoreApi_Clone(t *testing.T) {
	baseTempPath := "./.tmp"
	os.RemoveAll(baseTempPath)
	defer os.RemoveAll(baseTempPath)

	//given
	tests := map[string]struct {
		InputGitURL string
		OutputMeta  *icode.Meta
		OutputErr   error
	}{
		"success": {
			InputGitURL: "git@github.com:it-chain/tesseract.git",
			OutputMeta:  &icode.Meta{RepositoryName: "tesseract", GitUrl: "git@github.com:it-chain/tesseract.git", Path: baseTempPath + "/" + "tesseract"},
			OutputErr:   nil,
		},
		"fail": {
			InputGitURL: "git@github.com:nonono",
			OutputMeta:  nil,
			OutputErr:   errors.New("repository not found"),
		},
	}

	icodeApi, err := api.NewICodeGitStoreApi("testid", "testpw")
	assert.NoError(t, err)

	for testName, test := range tests {
		t.Logf("Running %s test, case: %s", t.Name(), testName)
		//when
		meta, err := icodeApi.Clone(baseTempPath, test.InputGitURL)
		if meta != nil {
			// icode ID 는 랜덤이기때문에 실데이터에서 주입
			// commit hash 는 repo 상황에따라 바뀌기 때문에 주입
			test.OutputMeta.ICodeID = meta.ICodeID
			test.OutputMeta.CommitHash = meta.CommitHash
		}

		//then
		assert.Equal(t, test.OutputMeta, meta)
		assert.Equal(t, test.OutputErr, err)

	}

}

func TestNewICodeGitStoreApi(t *testing.T) {
	id := "testid"
	pw := "testpw"
	icodeApi, err := api.NewICodeGitStoreApi(id, pw)
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(authUserID),
		Password: strings.TrimSpace(authUserPW),
	}
}

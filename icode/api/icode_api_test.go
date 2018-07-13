package api_test

import (
	"os"
	"testing"

	"path/filepath"

	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/it-chain/it-chain-Engine/icode/api"
	api2 "github.com/it-chain/it-chain-Engine/icode/infra/api"
	"github.com/it-chain/it-chain-Engine/icode/infra/service"
	"github.com/it-chain/tesseract"
	"github.com/stretchr/testify/assert"
)

func TestICodeApi_Deploy(t *testing.T) {
	//set data
	baseSaveUrl := "./.tmp"
	baseSaveUrl, err := filepath.Abs(baseSaveUrl)
	assert.NoError(t, err)
	GOPATH := os.Getenv("GOPATH")
	icodeGitUrl := "git@github.com:hea9549/test_icode"
	backupGitId := "validId"
	backupGitPw := "validPw"
	shPath := GOPATH + "/src/github.com/it-chain/tesseract/sh/default_setup.sh"
	//set mock repo
	mockRepo := GetMockRepo("deploy", nil)

	tesseractConfig := tesseract.Config{ShPath: shPath}
	containerService := service.NewTesseractContainerService(tesseractConfig, mockRepo)
	storeApi, err := api2.NewICodeGitStoreApi(backupGitId, backupGitPw)
	assert.NoError(t, err, "err in newIcodeGitStoreApi")
	icodeApi := api.NewIcodeApi(containerService, storeApi, mockRepo)
	err = icodeApi.Deploy(baseSaveUrl, icodeGitUrl)
	assert.NoError(t, err, "err in deploy")

	// todo event handler가없어서 icode 정보를 저장하지못하고잇음. 따라서 undeploy가 불가능 수동 삭제해줘야함
}

func TestICodeApi_UnDeploy(t *testing.T) {
	//set data
	baseSaveUrl := "./.tmp"
	baseSaveUrl, err := filepath.Abs(baseSaveUrl)
	assert.NoError(t, err)
	GOPATH := os.Getenv("GOPATH")
	icodeGitUrl := "git@github.com:hea9549/test_icode"
	backupGitId := "validId"
	backupGitPw := "validPw"
	shPath := GOPATH + "/src/github.com/it-chain/tesseract/sh/default_setup.sh"
	//set mock repo
	mockMeta := &icode.Meta{
		RepositoryName: "test_icode",
		GitUrl:         "git@github.com:hea9549/test_icode",
		Path:           filepath.Join(baseSaveUrl, "test_icode"),
	}
	mockRepo := GetMockRepo("unDeploy", mockMeta)

	tesseractConfig := tesseract.Config{ShPath: shPath}
	containerService := service.NewTesseractContainerService(tesseractConfig, mockRepo)
	storeApi, err := api2.NewICodeGitStoreApi(backupGitId, backupGitPw)
	assert.NoError(t, err, "err in newIcodeGitStoreApi")
	icodeApi := api.NewIcodeApi(containerService, storeApi, mockRepo)
	meta, err := icodeApi.Deploy(baseSaveUrl, icodeGitUrl)
	assert.NoError(t, err, "err in deploy")
	mockMeta = meta
	err = icodeApi.UnDeploy(mockMeta.ICodeID)
	assert.NoError(t, err, "err in unDeploy")

}

func GetMockRepo(testName string, mockData *icode.Meta) icode.ReadOnlyMetaRepository {
	if testName == "deploy" {
		return &mockRepo{}
	} else if testName == "unDeploy" {
		return &mockRepo{
			MockMeta: mockData,
		}
	}

}

type mockRepo struct {
	MockMeta *icode.Meta
}

func (m *mockRepo) FindById(id icode.ID) (*icode.Meta, error) {
	return m.MockMeta, nil
}

func (m *mockRepo) FindByGitURL(url string) (*icode.Meta, error) {
	return m.MockMeta, nil
}

func (m *mockRepo) FindAll() ([]*icode.Meta, error) {
	list := make([]*icode.Meta, 0)
	list = append(list, m.MockMeta)
	return list, nil
}

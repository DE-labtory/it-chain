package api_test

import (
	"os"
	"testing"

	"path/filepath"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/it-chain/it-chain-Engine/icode/api"
	api2 "github.com/it-chain/it-chain-Engine/icode/infra/api"
	"github.com/it-chain/it-chain-Engine/icode/infra/service"
	"github.com/it-chain/midgard"
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
	_, err = icodeApi.Deploy(baseSaveUrl, icodeGitUrl)
	assert.NoError(t, err, "err in deploy")
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
	eventstore.InitForMock(mockEventStore{})
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

func TestICodeApi_Invoke(t *testing.T) {
	baseSaveUrl := "./.tmp"
	baseSaveUrl, err := filepath.Abs(baseSaveUrl)
	assert.NoError(t, err)
	GOPATH := os.Getenv("GOPATH")
	icodeGitUrl := "git@github.com:hea9549/test_icode"
	backupGitId := "validId"
	backupGitPw := "validPw"
	shPath := GOPATH + "/src/github.com/it-chain/tesseract/sh/default_setup.sh"
	eventstore.InitForMock(mockEventStore{})
	mockMeta := &icode.Meta{
		RepositoryName: "test_icode",
		GitUrl:         "git@github.com:hea9549/test_icode",
		Path:           filepath.Join(baseSaveUrl, "test_icode"),
	}

	// txInvoke data 준비
	txInvoke := icode.Transaction{
		TxId: "1",
		TxData: icode.TxData{
			Jsonrpc: "2.0",
			Method:  "invoke",
			Params: icode.Param{
				Function: "initA",
				Args:     nil,
			},
		},
	}
	// txQuery data 준비
	txQuery := icode.Transaction{
		TxId: "2",
		TxData: icode.TxData{
			Jsonrpc: "2.0",
			Method:  "query",
			Params: icode.Param{
				Function: "getA",
				Args:     nil,
			},
		},
	}
	// mock repo 생성
	mockRepo := GetMockRepo("invoke", mockMeta)

	//teseeract 설정
	tesseractConfig := tesseract.Config{ShPath: shPath}
	containerService := service.NewTesseractContainerService(tesseractConfig, mockRepo)

	//storeApi 설정
	storeApi, err := api2.NewICodeGitStoreApi(backupGitId, backupGitPw)
	assert.NoError(t, err, "err in newIcodeGitStoreApi")

	//icodeApi 설정
	icodeApi := api.NewIcodeApi(containerService, storeApi, mockRepo)

	// deploy 시도
	meta, err := icodeApi.Deploy(baseSaveUrl, icodeGitUrl)
	assert.NoError(t, err, "err in deploy")

	// icode 정보 주입
	txInvoke.TxData.ID = meta.ICodeID
	txInvoke.TxData.ICodeID = meta.ICodeID
	txQuery.TxData.ID = meta.ICodeID
	txQuery.TxData.ICodeID = meta.ICodeID
	// Txs 데이터 준비
	invokeTxs := make([]icode.Transaction, 0)
	invokeTxs = append(invokeTxs, txInvoke)

	// invoke 시도
	invokeResults := icodeApi.Invoke(invokeTxs)

	// 결과 확인
	assert.Equal(t, true, invokeResults[0].Success)

	// query 시도
	queryResult, err := icodeApi.Query(txQuery)

	// 결과 확인
	assert.NoError(t, err, "err in query")
	assert.Equal(t, "0", queryResult.Data["A"], "diff in A data")

	// docker close
	err = icodeApi.UnDeploy(meta.ICodeID)
	assert.NoError(t, err)
}

func GetMockRepo(testName string, mockData *icode.Meta) icode.ReadOnlyMetaRepository {
	if testName == "deploy" {
		return &mockRepo{
			MockMeta: &icode.Meta{},
		}
	} else {
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
	return list, nil
}

type mockEventStore struct {
}

func (m mockEventStore) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return nil
}
func (m mockEventStore) Save(aggregateID string, events ...midgard.Event) error {
	return nil
}
func (m mockEventStore) Close() {

}

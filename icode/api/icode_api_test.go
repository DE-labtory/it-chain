package api_test

import (
	"github.com/it-chain/midgard"
)
/*
func TestICodeApi_Deploy(t *testing.T) {
	//set data
	baseSaveUrl := "./.tmp"
	sshPath := ""
	baseSaveUrl, err := filepath.Abs(baseSaveUrl)
	assert.NoError(t, err)
	GOPATH := os.Getenv("GOPATH")
	icodeGitUrl := "git@github.com:hea9549/test_icode"
	backupGitId := "validId"
	backupGitPw := "validPw"
	shPath := GOPATH + "/src/github.com/it-chain/tesseract/sh/default_setup.sh"
	//set mock repo

	tesseractConfig := tesseract.Config{ShPath: shPath}
	containerService := service.NewTesseractContainerService(tesseractConfig)
	storeApi, err := api2.NewICodeGitStoreApi(backupGitId, backupGitPw)
	assert.NoError(t, err, "err in newIcodeGitStoreApi")
	icodeApi := api.NewIcodeApi(containerService, storeApi)
	_, err = icodeApi.Deploy("1", baseSaveUrl, icodeGitUrl, sshPath)
	assert.NoError(t, err, "err in deploy")
}
*/
/*func TestICodeApi_UnDeploy(t *testing.T) {
	//set data
	baseSaveUrl := "./.tmp"
	baseSaveUrl, err := filepath.Abs(baseSaveUrl)
	sshPath := ""
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
	tesseractConfig := tesseract.Config{ShPath: shPath}
	containerService := service.NewTesseractContainerService(tesseractConfig)
	storeApi, err := api2.NewICodeGitStoreApi(backupGitId, backupGitPw)
	assert.NoError(t, err, "err in newIcodeGitStoreApi")
	icodeApi := api.NewIcodeApi(containerService, storeApi)
	meta, err := icodeApi.Deploy("1", baseSaveUrl, icodeGitUrl, sshPath)
	assert.NoError(t, err, "err in deploy")
	mockMeta = meta
	err = icodeApi.UnDeploy(mockMeta.ICodeID)
	assert.NoError(t, err, "err in unDeploy")

}*/
/*
func TestICodeApi_Invoke(t *testing.T) {
	baseSaveUrl := "./.tmp"
	baseSaveUrl, err := filepath.Abs(baseSaveUrl)
	sshPath := ""
	assert.NoError(t, err)
	GOPATH := os.Getenv("GOPATH")
	icodeGitUrl := "git@github.com:hea9549/test_icode"
	backupGitId := "validId"
	backupGitPw := "validPw"
	shPath := GOPATH + "/src/github.com/it-chain/tesseract/sh/default_setup.sh"
	eventstore.InitForMock(mockEventStore{})

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

	//teseeract 설정
	tesseractConfig := tesseract.Config{ShPath: shPath}
	containerService := service.NewTesseractContainerService(tesseractConfig)

	//storeApi 설정
	storeApi, err := api2.NewICodeGitStoreApi(backupGitId, backupGitPw)
	assert.NoError(t, err, "err in newIcodeGitStoreApi")

	//icodeApi 설정
	icodeApi := api.NewIcodeApi(containerService, storeApi)

	// deploy 시도
	meta, err := icodeApi.Deploy("1", baseSaveUrl, icodeGitUrl, sshPath)
	assert.NoError(t, err, "err in deploy")

	// icode 정보 주입
	txInvoke.TxData.ID = meta.ICodeID
	txInvoke.TxData.ICodeID = meta.ICodeID
	txQuery.TxData.ID = meta.ICodeID
	txQuery.TxData.ICodeID = meta.ICodeID
	// Txs 데이터 준비

	// invoke 시도
	invokeResults := icodeApi.Invoke(txInvoke)

	// 결과 확인
	assert.Equal(t, true, invokeResults.Success)

	// query 시도
	queryResult := icodeApi.Query(txQuery)

	// 결과 확인
	assert.Equal(t, "0", queryResult.Data["A"], "diff in A data")

	// docker close
	err = icodeApi.UnDeploy(meta.ICodeID)
	assert.NoError(t, err)
}*/

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

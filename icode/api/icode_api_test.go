/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	icodeApi := git.NewIcodeApi(containerService, storeApi)
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
	icodeApi := git.NewIcodeApi(containerService, storeApi)
	meta, err := icodeApi.Deploy("1", baseSaveUrl, icodeGitUrl, sshPath)
	assert.NoError(t, err, "err in deploy")
	mockMeta = meta
	err = icodeApi.UnDeploy(mockMeta.ICodeID)
	assert.NoError(t, err, "err in unDeploy")

}*/

//func TestICodeApi_Invoke(t *testing.T) {
//
//	baseSaveUrl := "./.tmp"
//	baseSaveUrl, err := filepath.Abs(baseSaveUrl)
//	sshPath := "../../id_rsa"
//	assert.NoError(t, err)
//	GOPATH := os.Getenv("GOPATH")
//	icodeGitUrl := "git@github.com:junbeomlee/test_icode"
//	shPath := GOPATH + "/src/github.com/it-chain/tesseract/sh/default_setup.sh"
//	eventstore.InitForMock(mockEventStore{})
//
//	// txInvoke data 준비
//	txInvoke := icode.Transaction{
//		TxId:      "1",
//		TimeStamp: time.Now(),
//		Jsonrpc:   "2.0",
//		Method:    "invoke",
//		ICodeID:   "",
//		Function:  "initA",
//		Args:      nil,
//	}
//	// txQuery data 준비
//	txQuery := icode.Transaction{
//		TxId:      "2",
//		TimeStamp: time.Now(),
//		Jsonrpc:   "2.0",
//		Method:    "query",
//		ICodeID:   "",
//		Function:  "getA",
//		Args:      nil,
//	}
//	// mock repo 생성
//
//	//teseeract 설정
//	tesseractConfig := tesseract.Config{ShPath: shPath}
//	containerService := service.NewTesseractContainerService(tesseractConfig)
//
//	//storeApi 설정
//	storeApi := git.NewRepositoryService()
//	assert.NoError(t, err, "err in newIcodeGitStoreApi")
//
//	//icodeApi 설정
//	icodeApi := api.NewIcodeApi(containerService, storeApi)
//
//	// deploy 시도
//	meta, err := icodeApi.Deploy("1", baseSaveUrl, icodeGitUrl, sshPath)
//	assert.NoError(t, err, "err in deploy")
//
//	// icode 정보 주입
//	txInvoke.ICodeID = meta.ICodeID
//	txQuery.ICodeID = meta.ICodeID
//	// Txs 데이터 준비
//
//	// invoke 시도
//	invokeResults := icodeApi.Invoke(txInvoke)
//
//	// 결과 확인
//	assert.Equal(t, true, invokeResults.Success)
//
//	fmt.Println(invokeResults)
//
//	// query 시도
//	queryResult := icodeApi.Query(txQuery.ICodeID, txQuery.Function, txQuery.Args)
//
//	// 결과 확인
//	assert.Equal(t, "0", queryResult.Data["A"], "diff in A data")
//
//	// docker close
//	err = icodeApi.UnDeploy(meta.ICodeID)
//	assert.NoError(t, err)
//}

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

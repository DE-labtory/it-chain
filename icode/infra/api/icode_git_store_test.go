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

/*
func TestICodeGitStoreApi_Clone(t *testing.T) {
	baseTempPath := "./.tmp"
	sshPath := ""
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
		meta, err := icodeApi.Clone(testName, baseTempPath, test.InputGitURL, sshPath)
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
	type Input struct {
		Id string
		Pw string
	}

	//given
	tests := map[string]struct {
		Input  Input
		Output string
	}{
		"validAccountCase": {
			Input: Input{
				// for test, write valid github id,pw
				Id: "validId",
				Pw: "validPw",
			},
			Output: "",
		},
		"invalidAccountCase": {
			Input: Input{
				Id: "invalidId",
				Pw: "invalidPw",
			},
			Output: "GET https://api.github.com/user: 401 Bad credentials []",
		},
	}

	for testName, test := range tests {
		t.Logf("Running %s test, case: %s", t.Name(), testName)
		//when
		_, err := api.NewICodeGitStoreApi(test.Input.Id, test.Input.Pw)

		if err != nil {
			assert.Equal(t, err.Error(), test.Output)
		} else {
			assert.Equal(t, "", test.Output)
		}
	}
}

//todo push를 어떻게 확인할지. 단순 레포 리스트만확인하면 createRepo테스트임. push 함수의 err만 체크해도되는지?
func TestICodeGitStoreApi_Push(t *testing.T) {
	validId := "validId"
	validPw := "validPw"
	baseTempPath := "./.tmp"
	sshPath := ""
	os.RemoveAll(baseTempPath)
	defer os.RemoveAll(baseTempPath)

	storeApi, err := api.NewICodeGitStoreApi(validId, validPw)

	assert.NoError(t, err)
	meta, err := storeApi.Clone("1", baseTempPath, "git@github.com:it-chain/heimdall.git", sshPath)
	err = storeApi.Push(*meta)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = storeApi.AuthClient.Repositories.Delete(ctx, validId, meta.RepositoryName)
	assert.NoError(t, err)
}
*/

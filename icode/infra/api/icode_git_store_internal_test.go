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

package api

/*
func TestCreateRepository(t *testing.T) {
	//given
	validId := "validId"
	validPw := "validPw"
	storeApi, err := NewICodeGitStoreApi(validId, validPw)
	assert.NoError(t, err, "icode store api 생성 실패")
	tests := map[string]struct {
		Input  string
		Output error
	}{
		"success": {
			Input:  "testingRepo",
			Output: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running '%s' test, caseName: %s", t.Name(), testName)

		//given
		err = storeApi.createRepository(test.Input)
		assert.NoError(t, err, "레포생성 실패")

		//then
		repos, _, err := storeApi.AuthClient.Repositories.List(context.Background(), validId, nil)
		assert.NoError(t, err, "레포 목록 불러오기 실패")
		found := false
		for _, repo := range repos {
			if repo.GetName() == test.Input {
				found = true
				break
			}
		}
		assert.Equalf(t, true, found, "cant find created repo in ID : %s", validId)
		assert.NoError(t, err)

		//after
		//생성확인후 삭제로직
		ctx := context.Background()
		_, err = storeApi.AuthClient.Repositories.Delete(ctx, validId, test.Input)
		assert.NoError(t, err, "레포생성후 삭제실패")
	}
}

func TestChangeRemote(t *testing.T) {
	panic("impl plz")
}

func TestNameFromGitUrl(t *testing.T) {
	panic("impl plz")
}
*/

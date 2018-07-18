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
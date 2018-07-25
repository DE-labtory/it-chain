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

package git_test

import (
	"errors"
	"os"
	"testing"

	"github.com/it-chain/engine/icode"
	"github.com/it-chain/engine/icode/infra/git"
	"github.com/stretchr/testify/assert"
)

func TestICodeGitStoreApi_Clone(t *testing.T) {
	baseTempPath := "./.tmp"
	sshPath := "./id_rsa"
	err, teardown := generatePriKey(sshPath)
	assert.NoError(t, err, "err in generate pri key file")
	defer teardown()
	os.RemoveAll(baseTempPath)
	defer os.RemoveAll(baseTempPath)

	//given
	tests := map[string]struct {
		InputGitURL string
		OutputMeta  *icode.Meta
		OutputErr   error
	}{
		"success": {
			InputGitURL: "git@github.com:hea9549/test_icode.git",
			OutputMeta:  &icode.Meta{RepositoryName: "test_icode", GitUrl: "git@github.com:hea9549/test_icode.git", Path: baseTempPath + "/" + "test_icode"},
			OutputErr:   nil,
		},
		"fail": {
			InputGitURL: "git@github.com:nonono",
			OutputMeta:  nil,
			OutputErr:   errors.New("repository not found"),
		},
	}

	icodeApi := git.NewRepositoryService()

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

func generatePriKey(path string) (error, func() error) {
	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}

	file, err := os.Create(path)
	if err != nil {
		return err, nil
	}

	key := "-----BEGIN RSA PRIVATE KEY-----\n" +
		"MIIEpAIBAAKCAQEAsTb2shAhgPcMKDqrmqjU4lzLsIWuzBuaEUQ1gVqNlvLr2Pks\n" +
		"rKvsRLOyYWhLgq405p8yHzwtaKGvq79rCEC3eLd906vA4WKpScqNJbEguCWapIvC\n" +
		"PGXCKYkmNO3AwNvbcEKpYjhI2Ypst7nr04fKqHF0yQzyf+uh27099wEM9OOHaqTh\n" +
		"hRVfN4gWT5SWgoi0V884OpMuG7pYt8jUFZKnU5DpKwE/hIMcrYzjvfv+epnhiIC3\n" +
		"BGs/zzeXsPldBhXceqfdoDsSjDZqx0dYcWiQSjLwwT1MF2r51MDzfMmoABxuMgjY\n" +
		"jExK2a0IUYy+ngmR1G2XtCfRaVe0wr6XezDmMwIDAQABAoIBADaEIxYaEkR7O3kw\n" +
		"u1PDtmHAjETMiz5tC1NeeVtGwSH7rwQ7ezvPU8q6wRhoHjqgXtPHi4LCX3G9s64R\n" +
		"H9sVFZwETqgMQTTUxiFWN1+uAtPDdbRC7kjoQPfIIkHMFiz+NZ5uU29Mw1Rw2gsX\n" +
		"He4f6v8wj+29lug1U8CmkeZno1W+GeV5I68bzaxRY8d8l9Xgkz+EGyvlCy26hXnz\n" +
		"ko2ghayWo+hbC10iJ2viBcLp077EEhbhXJnUx7eO4KGgAk+NUbz63Il/eg23yW42\n" +
		"Lwlj+DImdpT5AgAL7HLRLhGWVisbDtvuS7zqJ+tTcx4fnYyDL9NddaRWFV/ScWRU\n" +
		"7e9mYCECgYEA3c/WFvEkDj3grs/oddN7aUn60oBzuLApqDPDOTNsa6WvW5dSvvVO\n" +
		"xacrRxjviwxktMe29C7A9TZzw3UGm050g0BYBV+aom7/qBNK24MxCipIof5u/0jk\n" +
		"2nZdHYlK1fyCg3FJSQRxe35S1rQRHNq4h8YZnas6IwaXzsPWJyEodBECgYEAzIdu\n" +
		"NcG5gyfCPwUOz407hjL6mKUn9xsDpxND1JegBsQi47f8N1d64bwreRq8plDQjzZZ\n" +
		"5nMP5GATEtCulUfX+8FvICFJ1QIbLdjT05kqxJVjBCSsQfMOcIyGl02hvgwrW5xn\n" +
		"vzvL3ecCs1DUboAZDOstTJX0GOVcf1mHKGOg6gMCgYEAuEcGF1NJYCeaNcF24ATN\n" +
		"v8B9iEq9WU/Jm/s9EpWdWqVw1UgXr5v/UIg8lTmrMTsfo21UmvMIze/qJxVfYsHA\n" +
		"XJalSfmOb6qF7W3xwALzR/2vEB5guuglcTXq0DISoUrCZ09D0kzFtxJQ4h0BJpaz\n" +
		"veEfwPTTPOwqTNY4YZPTlvECgYAcuIvmapzVaRji5p/sz2Vjc/cyxkZ1ccqyhIcK\n" +
		"7HvhV1ua5LQ7RUKRPm5QZEvHgyO2aKh5LwE1TbR/+OP7PIp85O3o8iO/ELumVYNx\n" +
		"fFnAH0Y3R7sUy7/kWCdyScmDuYvBIpaDCS+Yqp77dUdPeReLc975mTkc4eB6VaUg\n" +
		"K9Vl5wKBgQCwThsDhlO4cHR+Ai7ZCe3pfq4t1z01XxNWLAX35T1TkMX4V51fi/MY\n" +
		"whsrbwTgsyW3Pdma1z5D+39wUYzZ2d+40fB0LPT4Q8b4pSwH7yjt294gwFIQvUlA\n" +
		"ylf3G/CjUTPvoz2HryxqoSjcYptaHG2F4C767Rx+V6pV4AP1hYp0dg==\n" +
		"-----END RSA PRIVATE KEY-----\n"
	file.WriteString(key)
	return nil, func() error {
		return file.Close()
	}
}

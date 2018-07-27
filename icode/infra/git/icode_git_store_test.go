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
		"MIIEowIBAAKCAQEAn7DJ1LzFa09rLOJAZlMWXu2/TRXe24b+73FLXDl0g/3clnqX\n" +
		"snuPWwLlUJ9YPr5EuVTAK09kStuhU+oxr9sxcyWlaQE+QX3uweYiKELK+pPRb3Bl\n" +
		"D9PmXXeU/tbFa+ey2sDMus/VWG1I+xG/Zwp13sY5pXK0l9IgMkg8TYj1jHsuxjhP\n" +
		"hvVi4a7PqnAfAW256zKiM+/3ELWTQZu4JuGKGJ2REJZMTCMz5d5Q/2H2jqPDMymM\n" +
		"23xyfFVFxHoce6W2+0sJjwSU5sxP6FgZxDR2o5EW/Yhk48/3fty5CrMOueaHtytn\n" +
		"dhYineLIXcz24VM6JKkj0zf6rmctRrXTDFXwcQIBJQKCAQB9KbLy2SUcwbTd/XC+\n" +
		"5y01pZXwT2l7oSit1Va959fdExStS1RNnx1VK8h9dfkOlSEbo0qCz4X1e7XMJkKe\n" +
		"lwsVD6soyZ9/cIqzpol5gHWT91Ef/iWydaa39fFSHfSuhStRgltGT96RammiMIFC\n" +
		"7IXl9U/G3V0dDHoZkpAhGFHBIQcwhhQBlDSLy4YmdQBCAey0hNkztYOfrUraeS61\n" +
		"ojvKsKEXBdahWw3NPgLjtctDo5xPpACFspPc2m7GrUs2W7HF8GeTcp3WtEJP1Mpc\n" +
		"hYey3RTzz4aIuTiuNHUWoKnpuF92OdCp9aXmROjoLIiiTagNgi8Tm0A+fJM9LjV3\n" +
		"1fedAoGBAOieGPUZYARlziO+PUxAY5sE2UDNj10m1PUH+Q91wI7OtLAFw7aBBwXY\n" +
		"zo+iZMpVNTv8qB0052yX/AzotWxSk++ky7lH3v6l7I9d0axEjuyZO6K5vm3JTyyI\n" +
		"Yv2ZokHaDoK70HsjIAEr2YyPzqQhvdTsM+Sfi2kjhTbrx+Vq1ULNAoGBAK++Ey8/\n" +
		"blk8Ru8cPHFSVx0I3AWo0UpBdkeM0xTSN4LdX8zknv+Ke96mGOnQh/y/0r/IEUus\n" +
		"jd9O++lRCuZE0cE2pN3I0SUmHROpFdU+kT4C0XQD5IBsEA/jP6GjpOtC4/Ud5QVo\n" +
		"dhjQqHyUHce7ipeizyZy71E7KInyY4YaQ4w1AoGBAIpQKoPlkwmIpBVArtNkjkB/\n" +
		"bGu4fsHED+uj4DK0t67bxWG+PQS7a/Wjgb0wIED0ZNePTzP18WMp79BTBBbk/gQc\n" +
		"zCj6TT13ahAORVGsOU5o8wbPTqIkq53wOtv6fCcntZnXdQrySmGPscky6ZIGOYWF\n" +
		"hqOdInXdxNSMMasBQIiVAoGBAKF+OyR/eih8lDWrPnX4oxPC87IsbUsnZcU1TFhS\n" +
		"eDMQnTjSFZDgHsyYoUWPMNprFLA7Tinc1WVrLKzisBHeX+H63LAgUXwVL3nZpVwr\n" +
		"qA94NheU38/0mSNbM42dS3Bm+vzrCc2eQwIL2RiIG12XlByjcj8B1P6JY4WuTZ3S\n" +
		"8fZ1AoGBANJI5iesUxyyDOcdujxkQ1aGguF9DMuauOfkNjSRQM5MFISh99Gtb+9F\n" +
		"JUNN2VBPEKKVCKVkXw+a1/DwwrJjDBy3Jld+b7T8g4MVzfMOMr/iQxBEUlF1ym/f\n" +
		"x35olce90ve7X7U1+FcG8YYlyvE7SCS1+TRSECnQ/FD1WpcBWhnR\n" +
		"-----END RSA PRIVATE KEY-----\n"
	file.WriteString(key)
	return nil, func() error {
		return file.Close()
	}
}

## Contribution 방법

0.  https://github.com/it-chain/engine 을 `fork`하기
1.  작업 할 내용 이슈에 등록하기 or 이슈에 이미 등록된 작업 선택 (작업중인 이슈에 댓글로 표시)
2.  이슈로 등록된 작업 수행
3.  수행한 작업의 문서, 단위 테스트 필수로 작성하기
4.  it-chain/engine 폴더에서 `go test -v ./... `로 테스트 수행, 
5.  goimports로 코드 포맷팅(issue #108) ```goimports -w ./```
6.  모든 테스트 통과시 작업한 이슈를 레퍼런스하여 `develop` 브랜치에 Pull Request
7.  travis ci 빌드 pass와 1명 이상의 approve를 받으면 `develop` 브랜치에 merge
8.  `master` 는 모든 테스트 케이스를 통과하며 빌드 에러가 없고 milestone 지점에 merge

### 브랜치 관리 규칙

* `master` : 릴리즈 수준의 코드만 merge.
* `develop` : 개발중인 테스트 완료된 코드만 merge.

### CLA

To get started, <a href="https://www.clahub.com/agreements/it-chain/engine">sign the Contributor License Agreement</a>.

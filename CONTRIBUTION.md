## Contribution 방법

0.  작업 할 내용 이슈에 등록하기 
1.  기능을 추가, 삭제, 업데이트, 이슈 등록, 문서 수정, 추가, 삭제 후 이슈 레퍼런스 하여 Pull Request
2.  구현한 기능에 관하여 문서작성 필수
3.  작성한 함수의 모든 단위테스트 작성
4.  서비스별 브랜치에서 작업
    1.  포맷 브랜치 예시: `feature/consensus`, `feature/blockchain` 등등
    2.  서비스별 브랜치를 구성하는 것으로 추천
5.  개발 후 `develop` 브랜치에 Pull Request
6.  `master` 는 모든 테스트 케이스를 통과하며 빌드 에러가 없고 milestone 지점에 merge

### 브랜치 관리 규칙

* `master` : 릴리즈 수준의 코드만 merge.
* `develop` : 서비스 단위로 완료된 코드만 merge.
* `feature/{service-name}/{feature or issue-title}` : 실제로 작업하는 브랜치. 개인 PC 에서 작업 후 push 하는 브랜치로, 작성 중인 피처나 픽스 중인 이슈가 완료되면 이 브랜치에서 `feature/{service-name}`로 pull request. (e.g. `feature/consensus/issue123` or `feature/consensus/election`)

#### `fork` 해서 코드를 작성할 경우

`fork` 해서 코드를 작성했을 경우, `develop` 으로 pull request 하는 것을 원칙으로 함

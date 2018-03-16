## 스마트 컨트랙트 <a name="SmartContract"></a>

![smartContract-implementation-deploy](../images/smartContract-implementation-deploy.png)

Smart Contract는 깃 저장소에 저장되어 있으며, 스마트 컨트랙트 서비스에 의해 실행됩니다. 안정성 및 보안을 위해 Docker 기반의 가상환경에서 테스트한 후, 문제가 없을 경우 실제 데이터 베이스에 반영됩니다.

- Git

  각 Smart Contract는 Git Repository에 저장됩니다.

- Docker VM

  Smart Contract를 실행하는 가상환경입니다. Docker VM에 Smart Contract와 최신 상태의 데이터베이스가 복사되면, 가상환경에서 실행되어 검증과정을 거칩니다.

- Smart Contract Service

  Git과 Docker VM을 관리하는 서비스입니다. Git에 Smart Contract를 push/clone하면 해당 서비스가 최신 상태의 데이터베이스와 Smart Contract를 복사하여 Docker VM에서 실행합니다.
  ​

#### Smart Contract 배포 순서도

![smartContract-implementation-seq](../images/smartContract-implementation-seq.png)

배포된 유저의 저장소는 검증된 Smart Contract 저장소에 저장되고 관리되어집니다. (아래 그림 참조)

| 유저 <br />저장소 <br />경로      | Smart Contract <br />저장소 <br />주소    | Smart Contract 파일 경로                   |
| -------------------------------- | ---------------------------------------- | ----------------------------------------- |
| A/a                              | {authenticated_git_id}/A_a               | It-chain/SmartContracts/A_a/{commit_hash} |
| B/b                              | {authenticated_git_id}/B_b               | It-chain/SmartContracts/B_b/{commit_hash} |
| C/c                              | {authenticated_git_id}/C_c               | It-chain/SmartContracts/C_c/{commit_hash} |

### World State DB
World State DB는 모든 트랜잭션이 실행된 후 마지막 상태에 대해서 저장합니다. 스마트컨트랙트(SmartContract)가 실행될 때 World state DB를 복사하여 이용한다.

| DB 이름              | 키              | 값                      | 설명                                                         |
| -------------------- | --------------- | ---------------------- | ------------------------------------------------------------ |
| WorldStateDB         | UserDefined Key | UserDefined Value      | 스마트컨트랙트 실행결과에 대한 정보를 저장                      |
| WaitingTransactionDB | Transaction ID  | Serialized Transaction | 트랜잭션 저장                                                 | 

### Author

[@hackurity01](https://github.com/hackurity01)
[@codeblv](https://github.com/codeblv)
[@codeblv](https://github.com/codeblv)

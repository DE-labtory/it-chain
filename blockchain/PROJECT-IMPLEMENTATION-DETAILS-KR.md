# 블록체인 <a name="BlockChain"></a>

![blockchain-implemeneation-logical](../images/blockchain-implemeneation-logical.png)

- 블록체인

  블록체인은 [해쉬](https://en.wikipedia.org/wiki/Cryptographic_hash_function)로 연결된 지속적으로 늘어나는 블록들의 리스트다.그 해쉬는 이전 블록의 링크역할을 한다.

- 블록

  블록은 블록 헤더와 블록 데이터로 이루어져있다. 그리고 Ledger의 블록체인 구조를 위하여 다음 블록은 블록 헤더를 해싱한 값을 가지고 있다. 블록헤더는 이전 블록의 해쉬 값과 머클 트리 루트의 해쉬 값을 가지고 있다. 블록데이터는 트랜잭션 리스트를 가지며 트랜잭션의 위변조를 효율적으로 관리하기 위해 머클트리를 가지고 있다.

- 트랜잭션

  스마트 컨트랙트를 수행하기 위한 작은 장치이다. 트랜잭션은 실제로 트랜잭션을 실행시키는 피어(노드)의 아이디, 트랜잭션 해쉬 값, 계약 내용을 담고 있는 트랜잭션 데이터를 가지고있다.

- 머클 트리

  머클 트리는 이진 트리이고 각 리프노드들은 블록 안의 트랜잭션의 해쉬값을 가진다. 루트 노드는 전체 트랜잭션을 나타내는 해쉬 값을 가진다. 여기서 전체 트랜잭션을 나타내는 해쉬값은 리프 노드부터 각각 자식 노드들의 해쉬 값들을 해쉬한 값을 말한다. 머클 트리는 O(1)시간만에 트랜잭션 정보들이 바뀌었는지 머클 트리의 루트 해쉬값을 통해서 확인할 수 있다. 더불어, 머클 트리는 원장 안에 있는 트랜잭션들의 유효성을 효율적으로 검증할 수 있다. 왜냐하면 블록 헤더는 머클 트리의 루트 노드의 해쉬 값을 가지고 있고 다음 블록은 현재 블록의 해쉬 값을 블록 헤더에 가지고 있기 때문이다. 그리고 머클 트리는 머클 경로(트랜잭션의 루트 노드까지의 자식 노드)를 구할 수 있기 때문에 특정 트랜잭션의 유효성을 O(lgN)시간 만에 구할 수 있다는 장점이 있다.


  ![blockchain-implementation-merkletree](../images/blockchain-implementation-merkletree.png)

## Database <a name="DB"></a>
블록체인은 구성에 따라 여러 유형의 데이터베이스에 저장된다. 현재는 levelhelper와 filehelper의 기능이 추가되어져 있다. 기본 DB는 levelDB이다. 블록의 해시값이나 블록의 번호, 트랜잭션 ID를 가지고 블록을 검색할 수 있습니다. 또한 트랜잭션ID를 가지고 해당하는 트랜잭션을 검색 가능하다. 만약 다른 데이터베이스를 사용하길 원한다면 blockchainleveldb에서 구현하고 blockchain_db_interface를 수정하세요.

### Related config
데이터베이스의 설정은 config.yaml에 정의된다.

- Type

  데이터베이스의 유형. 현재는 levelDB와 파일에 대한 몇몇의 helper기능이 지원된다.

- Leveldb

  leveldb의 구성
  
  | 키           | 설명                                            |
  | ------------ | --------------------------------------------    |
  | default_path | leveldb의 경로를 설정하지 않는다면 이 경로로 설정  |
  
### LevelDB

블록들은 키-값 스토리지인 leveldb에 저장된다.

- Blocks  

  블록들은 JSON의 형태로 직렬화되어지고 leveldb에 저장된다. 키값으로 블록의 해쉬 값과 번호가 사용된다.
  마지막 블록와 확정되지 않은 블록도 복구를 위해서 저장된다.
  
- Transactions

  모든 트랜잭션들은 직렬화되어지고 leveldb에 저장됩니다. 기본적으로 모든 트랜잭션들은 블록와 함께 저장된다.
  인덱싱을 위해서, 트랜잭션이 속한 블록의 해쉬 값도 저장된다. 트랜잭션ID는 키로 이용된다.
  
| DB 이름            | 키             | 값                           | 설명                                     |
| ----------------- | -------------- | ---------------------------- | ---------------------------------------- |
| block_hash        | BlockHash      | Serialized Block             | 블록해시를 이용하여 블록 저장              |
| block_number      | BlockNumber    | Block Hash                   | 블록 번호를 이용하여 블록 저장             |
| transaction       | Transaction ID | Serialized Transaction       | 트랜잭션 저장                             |
| unconfirmed_block | BlockHash      | Serialized unconfirmed block | 확정되지 않은 블록 저장                    |
| util              | Predefined Key | Depends on Key               | 여러 용도를 위한 DB                        |  

- util DB

  Util DB는 편의를 위해서 여러 항목을 저장한다.
  
  1) Key : last_block, Value : Serialized last block
  2) Key : unconfirmed_block, Value : Serialized unconfirmed block
  3) Key : transaction ID, Value : Blockhash of block that transaction is stored

- Snapshot

  leveldb에 저장된 world state db를 복사하기 위해서 LevelDB snapshot이 추가된다.
  
### File
블록의 메타데이터는 leveldb 또는 다른 키-값 데이터베이스에 저장된다. 블록의 바디부분은 파일에 저장된다.

- Blocks

  블록의 메타데이터는 JSON에 직렬화되어 leveldb에 저장된다. 블록바디에 삽입되는 데이터는 파일에 기록된다.
  
- Transactions

  트랜잭션 데이터는 file에 저장된다. 검색을 위해서 트랜잭션 ID를 키로 사용하는 키-값 데이터베이스에 파일의 정보를 저장한다.

### Author

[@emperorhan](https://github.com/emperorhan)
[@zeroFruit](https://github.com/zeroFruit)

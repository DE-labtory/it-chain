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
  
### World State DB

World State DB는 모든 트랜잭션이 실행된 후 마지막 상태에 대해서 저장합니다. 스마트컨트랙트(SmartContract)가 실행될 때 World state DB를 복사하여 이용한다.

| DB 이름              | 키              | 값                      | 설명                                                         |
| -------------------- | --------------- | ---------------------- | ------------------------------------------------------------ |
| WorldStateDB         | UserDefined Key | UserDefined Value      | 스마트컨트랙트 실행결과에 대한 정보를 저장                      |
| WaitingTransactionDB | Transaction ID  | Serialized Transaction | 트랜잭션 저장                                                 | 

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

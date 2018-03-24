# BlockChain

BlockChain은 it-chain-Engine에서 사용되는 **ledger**의 Block을 생성, 조회 및 검증한다.



## Block
Block에는 BlockHeader와 BlockData가 있다.

- **BlockHeader**
  - Block Number
  - Previous Block Hash
  - Version
  - Merkle Tree Root Hash
  - Time Stamp
  - Block Height
  - Block Status
  - Block Hash

- **BlockData**
  - Merkle Tree
  - Transaction List
  - Transaction Index Map
  - Merkle Tree Height
  - Transaction Count

## Key methods of Block

Block의 주요 메소드

- CreateNewBlock
    - 새로운 블록을 생성
- PutTransaction
    - 트랜잭션 리스트를 블록에 append
- MakeMerkleTree
    - 현재 트랜잭션 리스트를 이용하여 머클트리 생성
- VerifyTx
    - 머클패스를 이용하여 트랜잭션 검증
- VerifyBlock
    - 모든 트랜잭션을 검증하여 블록이 올바른지 검증
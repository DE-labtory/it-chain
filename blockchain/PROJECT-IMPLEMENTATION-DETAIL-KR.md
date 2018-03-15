# 블록체인 <a name="BlockChain"></a>

![blockchain-implemeneation-logical](../../images/blockchain-implemeneation-logical.png)

- 블록체인(BlockChain)

  블록체인은 [해쉬](https://en.wikipedia.org/wiki/Cryptographic_hash_function)로 연결된 지속적으로 늘어나는 블록들의 리스트다.그 해쉬는 이전 블록의 링크역할을 한다.

- 블록(Block)

  블록은 블록 헤더와 블록 데이터로 이루어져있다. 그리고 Ledger의 블록체인 구조를 위하여 다음 블록은 블록 헤더를 해싱한 값을 가지고 있다. 블록헤더는 이전 블록의 해쉬 값과 머클 트리 루트의 해쉬 값을 가지고 있다. 블록데이터는 트랜잭션 리스트를 가지며 트랜잭션의 위변조를 효율적으로 관리하기 위해 머클트리를 가지고 있다.

- 트랜잭션(Transaction)

  스마트 컨트랙트를 수행하기 위한 작은 장치이다. 트랜잭션은 실제로 트랜잭션을 실행시키는 피어(노드)의 아이디, 트랜잭션 해쉬 값, 계약 내용을 담고 있는 트랜잭션 데이터를 가지고있다.

- 머클 트리(Merkle Tree)

  머클 트리는 이진 트리이고 각 리프노드들은 블록 안의 트랜잭션의 해쉬값을 가진다. 루트 노드는 전체 트랜잭션을 나타내는 해쉬 값을 가진다. 여기서 전체 트랜잭션을 나타내는 해쉬값은 리프 노드부터 각각 자식 노드들의 해쉬 값들을 해쉬한 값을 말한다. 머클 트리는 O(1)시간만에 트랜잭션 정보들이 바뀌었는지 머클 트리의 루트 해쉬값을 통해서 확인할 수 있다. 더불어, 머클 트리는 원장 안에 있는 트랜잭션들의 유효성을 효율적으로 검증할 수 있다. 왜냐하면 블록 헤더는 머클 트리의 루트 노드의 해쉬 값을 가지고 있고 다음 블록은 현재 블록의 해쉬 값을 블록 헤더에 가지고 있기 때문이다. 그리고 머클 트리는 머클 경로(트랜잭션의 루트 노드까지의 자식 노드)를 구할 수 있기 때문에 특정 트랜잭션의 유효성을 O(lgN)시간 만에 구할 수 있다는 장점이 있다.


  ![blockchain-implementation-merkletree](../../images/blockchain-implementation-merkletree.png)

### Author

[@emperorhan](https://github.com/emperorhan)
[@zeroFruit](https://github.com/zeroFruit)

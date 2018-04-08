# Consensus

합의(Consensus)는 it-chain-Engine의 P2P Network의 구성원들 사이에서 수행되며 리더노드에 의해 시작된다. Consensus Module은 리더에 의해 시작된 **Block의 저장 순서에 대한 합의**만을 수행한다.

블록에 포함된 내용이 이미 잘 검증된 트랜잭션들이고, 각 트랜잭션의 결과가 deterministic 하다면, Block의 저장 순서만 정확히 합의할 경우 모든 노드가 언제나 동일한 결과물(world state)을 갖게 된다.



## Consensus 알고리즘

It-chain-Engine는 PBFT를 Block 합의 알고리즘으로 사용한다.



## Primary Data of Consensus

Consensus Module에서 다루는 핵심 데이터

- Consensus in progress (합의 중인 합의)
- ConsensusMsg
  - PreprepareMsg
  - PrepareMsg
  - CommitMsg
- Paliament(P2P Network의 구성원)



## Consume Event

- Parliament변경
  - LeaderChangeEvent 
  -  PeerConnectEvent
  -  PeerDisConnectEvent
- 합의할 Block 준비완료
- ConsensusMessageArriveEvent(3 type)
  - PreprepareMsg
  - PrepareMsg 
  - CommitMsg 



## Publish Event

- BlockConfirmEvent
- ConsensusMessagePublishEvent
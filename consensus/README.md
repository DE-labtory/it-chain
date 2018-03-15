# Consensus

합의(Consensus)는 it-chain-Engine의 P2P Network의 구성원들 사이에서 수행되며 리더노드에 의해 실행된다. Consensus Module은 리더에 의해 시작된 **Block의 합의**만을 수행한다.



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

- Paliament변경
  - ChangeLeader 
  -  AddPeer
  -  DeletePeer
- 합의할 Block 준비완료
- ConsensusMessage
  - ReceivePreprepareMsg 
  - ReceivePrepareMsg 
  - ReceiveCommitMsg 



## Publish Event

- ConfirmBlock
- SendConsensusMessage
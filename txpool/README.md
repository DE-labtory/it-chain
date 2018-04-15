# TxPool

TxPool은 노드의 트랜잭션을 관리한다. TxPool은 **일정조건**을 만족하면 다음과 같은 작업을 진행한다.

- 자신이 리더 노드일 경우 모여있던 트랜잭션을 블록으로 만들 준비가 됬다는 이벤트를 보낸다.
- 자신이 리더노드가 아닐경우 모여있던 트랜잭션을 리더노드에게 보낼 준비가 됬다는 이벤트를 보낸다.

# Primary Data of TxPool

TxPool에서 다루는 데이터

- ㅁㄴㅇㄹ
- ㄴㅇㄹ

- ㄴㅇㄹ

# Consume Event

- transaction 추가/삭제
    - CreateTxEvent
    - DeleteTxEvent

- Parliament 변경
    - LeaderChangeEvent
    - PeerDisconnectEvent

# PublishEvent

- TxReadyToBlockEvent
- TxSendToLeaderEvent
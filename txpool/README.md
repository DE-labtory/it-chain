# TxPool

TxPool은 노드의 트랜잭션을 관리한다. TxPool은 **일정조건**을 만족하면 다음과 같은 작업을 진행한다.

- 리더일경우에는 블록생성 조건에 의해 블록만들 준비 이벤트 (TxReadyToBlockEvent)를 날리고,
- 리더가 아닐경우에는 트랜잭션 전송 조건에 의해 리더 노드에게 트렌잭션을 보내는 이벤트 (TxSendToLeaderEvent)를 날린다

# Primary Data of TxPool

TxPool에서 다루는 데이터

- transaction

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
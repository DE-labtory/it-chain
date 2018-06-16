# TxPool

TxPool은 노드의 트랜잭션을 관리한다. TxPool은 **일정조건**을 만족하면 다음과 같은 작업을 진행한다.

- 리더일경우에는 블록생성 조건에 의해 블록만드는 command(ProposeBlock)를 날리고,
- 리더가 아닐경우에는 트랜잭션 전송 조건에 의해 리더 노드에게 트렌잭션을 보내는 command(SendLeaderTransactions)를 날린다
- tx와 관련된 event(생성, 삭제)를 수신하고 해당 tx를 변경한다.
- leader와 관련된 event를 수신하고 leader 정보가 변경되면 TxPool에서도 그에 맞게 변경한다.

## API
## Message Dispatcher
### ProposeBlock(transactions []txpool.Transaction)
block을 만들기 위한 transactions들을 blockchain에게 넘겨준다.
### SendLeaderTransactions(transactions []*txpool.Transaction, leader txpool.Leader)
leader에게 transactions을 보내준다.

## Event Handler

### HandleTxCreatedEvent(txCreatedEvent txpool.TxCreatedEvent)
event를 받으면 해당 tx를 저장한다.

### HandleTxDeletedEvent(txDeletedEvent txpool.TxDeletedEvent)
event를 받으면 해당 tx를 삭제한다.

### HandleLeaderChangedEvent(leaderChangedEvent txpool.LeaderChangedEvent)
event를 받으면 해당 leader 정보로 업데이트한다.

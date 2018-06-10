# P2P Component
## 서론
Peer 서비스는 피어노드의 생성, 삭제, 열람, 리더 노드의 선정 및 변경 등의 역할을 수행합니다.

## API



## Message Dispatcher
### RequestLeaderInfo(peer p2p.Node)

### DeliverLeaderInfo(toPeer p2p.Node, leader p2p.Node)

### RequestNodeList(peer p2p.Node)

### DeliverNodeList(toNode p2p.Node, nodeList []p2p.Node)






## Message Protocol and Types of Messages

### LeaderInfoRequestProtocol
request of leader info

Message: `LeaderInfoRequestMessage`

### LeaderInfoDeliverProtocol
delivery of leader info

Message: `LeaderInfoDelivery`

### NodeListRequestProtocol
request of node list

Message: `NodeListRequestMessage`

### NodeListDeliverProtocol
send contents of node repository to specific node

Message: `NodeListDelivery`










---

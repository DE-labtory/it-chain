# P2P Component
## 서론
Peer 서비스는 피어노드의 생성, 삭제, 열람, 리더 노드의 선정 및 변경 등의 역할을 수행합니다.

## API

## Message Dispatcher

### RequestLeaderInfo(peer p2p.Node)
request of leader info

### DeliverLeaderInfo(toPeer p2p.Node, leader p2p.Node)
delivery of leader info

### RequestNodeList(peer p2p.Node)
request of node list

### DeliverNodeList(toNode p2p.Node, nodeList []p2p.Node)
delivery of node list

### DelverNode(toNode p2p.NodeId, node p2p.Node)
deliver of single node

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

### NodeDeliverProtocol
send single node to specific node

Message: `NodeDelivery`










---

 # Grpc Gateway

GrpcGateway는 노드끼리의 통신을 담당한다. GrpcGateway는 bifrost(https://github.com/it-chain/bifrost)를 사용하여 grpc기반 p2p network를 구성한다.

- 다른 노드와의 네트워크 연결을 수행 및 관리 한다.
- 다른 노드에게 데이터 전송을 수행한다.
- 다른 노드로 부터 데이터를 받아 관련된 데이터 Event를 발생 시킨다.



# Consume Event

- MessageDeliveryEvent

# PublishEvent

- PeerConnectEvent

- PeerDisconnectEvent

- CMessageArriveEvent

  … 
# Grpc Gateway

GrpcGateway는 노드끼리의 통신을 담당한다. GrpcGateway는 bifrost(https://github.com/it-chain/bifrost)를 사용하여 grpc기반 p2p network를 구성한다.

- 다른 노드와의 네트워크 연결을 수행 및 관리 한다.
- 다른 노드에게 데이터 전송을 수행한다.
- 다른 노드로 부터 데이터를 받아 관련된 데이터 Event를 발생 시킨다.

이를 위해서는 grpc client 및 server의 생성 및 관리, protocol buffer 관리와 구체적인 grpc 서비스를 구현해야 하는데 이 모든 작업은 **bifrost** 라이브러리를 통해 구현한다.

# Methods
## start()
in gateway.go

## stop()
in gateway.go

# Structures
## Server
in server.go

gRPC 의 서버의 모든 기능을 수행한다.
응답을 대기하고, 새로운 서버를 생성하고, rpc connection이 생겼을 때에 대한 handler 역할을 수행한다.

### Server Structure
### onConnection
### Listen(ip string)
### onError(error error)
### Stop()
# Consume Event

- MessageDeliveryEvent

# PublishEvent

- PeerConnectEvent

- PeerDisconnectEvent

- CMessageArriveEvent

  …

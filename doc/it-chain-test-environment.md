# IT CHAIN P2P TEST ENVIRONMENT

## 서문

it-chain 은 여러 peer 된 p2p 네트워크 시스템으로 많은 핵심적인 기능들은 피어간 통신을 필요로 한다. 때문에 전체 네트워크가 완벽하게 구축되지 않은 상황에서 각 컴포넌트의 핵심 기능들을 테스트하기가 매우 어려우며, 이를 해결하기 위해 각 peer를 단일 노드의 프로세스로 추상화하여 프로세스가 통신을 통해 inter-network function 들을 효과적으로 테스트함을 그 목적으로 한다.

## BASIC CONCEPT

다음은 it-chain 의 단일 노드 테스트 환경에 대한 architecture 이다.

![단일 노드 테스트 환경](./images/[test]mock_process_design.png)

테스트 시에 단일 노드는 MockProcess 라는 구현체를 통해 추상화 될 수 있다. it-chain 테스트 환경에서 rabbitmq는 동작하지 않음을 가정하므로 amqp 의 동작이 가상으로 구현되어야 한다. MockProcess 는 특정 노드가 받는 모든 event 및 command 를 evenListenr와 commandLister 라는 channel에 주입시켜 주며, 실제 event handler, command handler에서 event 와 command가 소비되어 처리된다. 또한 각 컴포넌트는 evnet 혹은 command를 publis하는 과정을 기존처럼 midgard 를 통해 rabbitma에 전달하는 것이 아니라 MockProcess 내의 commandListener 혹은 eventListener에 직접 전달한다.

특히, 다른 노드로 message 를 전달하는 특정한 command 인 grpcDeliverCommand 의 경우 gateway component 의 역할을 대행하기 위해 command handler에서 처리되며, 기존에 message.deliver command 를 통해 gateway 에서 grpc 로 grpcDeliverCommand 를 보내던 것을 직접 다른 mockProcess의 commandReceiver에 grpcReceiveCommand를 전달하는 것으로 대체하였다. grpcReceiverCommand의 경우 실제 handler에서 동일하게 처리된다.

---

### AUTHOR

[@fronalnh](http://github.com/frontalnh)

---

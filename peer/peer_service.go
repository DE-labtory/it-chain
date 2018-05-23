package peer

// 현재 본 파일이 peer의 대부분의 service를 수행하는 것으로 보입니다.
// 특히 message_produce와 관련된 다양한 함수를 처리하는것으로 보이고
// message produce의 특성상 infra 에 포함된 messaging 기능을 처리하기 위해 message produce의 interface만 구현해 놓았다.

type Publish func(topic string, data []byte) error // type 문을 정의함으로써 해당 함수의 원형을 간단히 표현
// topic과 byte 배열인 data를 받아 error를 반환

type MessageProducer interface {
	RequestLeaderInfo(peer Peer) error
	DeliverLeaderInfo(toPeer Peer, leader Peer) error
	LeaderUpdateEvent(leader Peer) error
}

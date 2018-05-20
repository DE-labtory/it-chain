package service

import "github.com/it-chain/it-chain-Engine/peer/domain/model"

// 현재 본 파일이 peer의 대부분의 service를 수행하는 것으로 보입니다.
// 특히 message_produce와 관련된 다양한 함수를 처리하는것으로 보이고
// message produce의 특성상 infra 에 포함된 messaging 기능을 처리하기 위해 message produce의 interface만 구현해 놓은 것으로 보입니다.

// 하지만 messaging 과 관련된 함수를 infra에서 구현한 뒤에 실제적인 로직이 현 파일에서 구현되는 것이 좀더 명확하지 않을까 싶습니다.
// peer와 관련된 기능중 infra와 독립적인 내용이 이 구간에서 모두 구현되는 것이 맞는데 모든 구현이 infra에서 이루어져서 service 가 제 역할을 하지 못하는 것 같습니다
// 가령, 만약 새롭게 필요한 기능 중 하나가 messaging과도 관련되지 않고, repository와도 관련되지 않은 특정 기능이 있다면 본 service 디렉토리에서 이루어져야 하나
// 현재의 message_produce와 같은 기능들은 단지 interface로만 구현이 되어있는데 반대 새롭게 생기는 기능은 interface만 구현되어 있어 혼선이 빚어질 것 같습니다.
// 즉, rabbitmq 를 다루는 내용의 함수는 infra에서 구현이 되어지되, 해당 함수를 호출하여 본 파일에서 구현이 완료되는 것이 개발자 입장에서 보다 직관적일 것 같습니다.

// 추가적으로 본 파일 명은 peerService가 되고 실제 구현도 해당 파일에서 infra의 함수를 호출하여 이루어 지도록 하는 것이 명확할 것 같습니다.
// ---
// by frontalnh(namhoon)

type Publish func(topic string, data []byte) error // type 문을 정의함으로써 해당 함수의 원형을 간단히 표현
																									 // topic과 byte 배열인 data를 받아 error를 반환

type MessageProducer interface {
	RequestLeaderInfo(peer model.Peer) error
	DeliverLeaderInfo(toPeer model.Peer, leader model.Peer) error
	LeaderUpdateEvent(leader model.Peer) error
}

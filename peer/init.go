package peer

// 임의로 명칭을 붙였습니다.
// 기존의 peer.go 는 모델명으로 활용해야 하기에 이런이름을 임의로 붙였습니다

import (
	"github.com/it-chain/it-chain-Engine/conf"
)

// initiallize peer
//
func init() {
	bootNodeIp := conf.GetConfiguration().Common.BootNodeIp
	myIp := conf.GetConfiguration().Common.NodeIp
	myPeer := NewPeer(myIp, peerId) // peerId는 어디서 할당?
	messageProducer := NewMessageProducer()
	repository := NewPeerRepository("path") //repository 객체 생성
	leaderSelectionApi := NewLeaderSelectionApi(repository, messageProducer, myPeer)

	// 해당 노드를 leveldb에 저장
	repository.save(myPeer)

	// 부트노드 ip가 나의 ip와 같은 경우 나 자신을 리더로 설정한다.
	// 현재의 경우 bootNodeIp와 myIp가 동일하므로 무조건 사용자 스스로를 리더 피어로 설정하고 있다.
	if bootNodeIp == myIp {

		// 사용자를 리더로 선언
		err := leaderSelectionApi.messageProducer.LeaderUpdateEvent(myPeer)

		if err != nil {
			//todo mq error handle 해야하나?
		}
	}
}

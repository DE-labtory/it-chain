package peer

// 임의로 명칭을 붙였습니다.

import (
	"github.com/it-chain/it-chain-Engine/conf"
)

// initiallize peer
//
func init() {
	bootNodeIp := conf.GetConfiguration().Common.BootNodeIp
	myIp := conf.GetConfiguration().Common.NodeIp
	myInfo := NewPeer(myIp, peerId) // peerId는 어디서 할당?
	messageProducer := NewMessageProducer
	leaderSelectionApi := NewLeaderSelectionApi(repo repository.Peer, messageProducer service.MessageProducer, myInfo)

	// 부트노드 ip가 나의 ip와 같은 경우 나 자신을 리더로 설정한다.
	// 현재의 경우 bootNodeIp와 myIp가 동일하므로 무조건 사용자 스스로를 리더 피어로 설정하고 있다.
	if bootNodeIp == myIp {

		// 사용자를 리더로 선언
		err := leaderSelectionApi.messageProducer.LeaderUpdateEvent(myInfo)

		if err != nil {
			//todo mq error handle 해야하나?
		}
	}
}

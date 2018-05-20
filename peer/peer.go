package peer

import "github.com/it-chain/it-chain-Engine/conf"

func init() {
	bootNodeIp := conf.GetConfiguration().Common.BootNodeIp
	myIp := conf.GetConfiguration().Common.NodeIp

	// 부트노드 ip가 나의 ip와 같은 경우 나 자신을 리더로 설정한다.
	// 현재의 경우 bootNodeIp와 myIp가 동일하므로 무조건 사용자 스스로를 리더 피어로 설정하고 있다.
	if bootNodeIp == myIp {
		err := leaderSelectionApi.messageProducer.LeaderUpdateEvent(*myInfo)
		if err != nil {
			//todo mq error handle 해야하나?
		}
	}
}

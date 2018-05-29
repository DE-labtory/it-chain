package p2p

import (
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
	"github.com/it-chain/it-chain-Engine/p2p/infra/messaging"
	"github.com/it-chain/midgard"
)


func init() {
	bootNodeIp := conf.GetConfiguration().Common.BootNodeIp
	myIp := conf.GetConfiguration().Common.NodeIp
	myNode := NewNode(myIp, nodeId) // nodeId는 어디서 할당?
	//messageDispatcher := messaging.NewMessageDispatcher(midgard.Publisher) midgard 주입부분 => midgard doc 완성 후로 보류
	repository := leveldb.NewNodeRepository("path") //repository 객체 생성
	leaderSelectionApi := NewLeaderSelectionApi(repository, messageDispatcher, myNode)

	// 해당 노드를 leveldb에 저장
	repository.save(myNode)

	// 부트노드 ip가 나의 ip와 같은 경우 나 자신을 리더로 설정한다.
	// 현재의 경우 bootNodeIp와 myIp가 동일하므로 무조건 사용자 스스로를 리더 피어로 설정하고 있다.
	if bootNodeIp == myIp {

		// 사용자를 리더로 선언
		err := leaderSelectionApi.messageDispatcher.LeaderUpdateEvent(myNode)

		if err != nil {
			//todo mq error handle 해야하나?
		}
	}
}

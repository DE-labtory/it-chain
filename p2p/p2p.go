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

}

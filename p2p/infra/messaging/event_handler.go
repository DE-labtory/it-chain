package messaging

import (
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
	"github.com/it-chain/midgard"
)

type NodeEventHandler struct {
	nodeRepository   			p2p.NodeRepository
	leaderRepository 			p2p.LeaderRepository
	messageDispatcher       	*MessageDispatcher
}

func NewNodeEventHandler(nodeRepo *leveldb.NodeRepository, leaderRepo *leveldb.LeaderRepository, messageDispatcher *MessageDispatcher) *NodeEventHandler {
	return &NodeEventHandler{
		nodeRepository:   			nodeRepo,
		leaderRepository: 			leaderRepo,
		messageDispatcher:     	  	messageDispatcher,
	}
}

//todo conn후 peer에 추가하고 peer 추가됬다는 이벤트를 날려줘야 하나 말아야하나 고민이슈
//p2p 에서 노드가 추가된 뒤에 노드가 추가되었다는 사실을 알리는 식으로 doc/ 에 시나리오 작성해 두었습니다.
//참고 바랍니다.
func (neh *NodeEventHandler) HandleConnCreatedEvent(event p2p.ConnectionCreatedEvent) {
	id := p2p.NodeId(event.ID) //id 생성은 어떻게 동작하는가? aggregateID를 사용?!
	address := event.Address
	node := p2p.NewNode(address, id)
	neh.nodeRepository.Save(*node)


	if neh.leaderRepository.GetLeader() == nil {
		err := neh.messageDispatcher.RequestLeaderInfo(*node)
		if err != nil {
			log.Println(err.Error())
		}
	}

	//set

	// 노드 생성 이벤트 날림
	err := neh.messageDispatcher.publisher.Publish("Event", "Connection", p2p.NodeCreatedEvent{
		EventModel: midgard.EventModel{
			ID: connection.GetID(),
		},
	})

	if err != nil{
		log.Println(err.Error())
	}

}

//todo conn disconnect event 구현
func (neh *NodeEventHandler) HandleConnDisconnectEvent(event p2p.ConnectionDisconnectedEvent) {
	panic("need to implement")
}

//save Leader when LeaderReceivedEvent Detected
func (neh *NodeEventHandler) HandleLeaderReceivedEvent(event p2p.LeaderReceivedEvent){
	panic("implement!")
}

//save all nodes when NodeListReceivedEvent Detected
func (neh *NodeEventHandler) HandleNodeListReceivedEvent(event p2p.NodeListReceivedEvent){
	panic("implement!")
}

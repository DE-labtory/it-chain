package messaging

import (
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
)

type NodeEventHandler struct {
	nodeRepository   p2p.NodeRepository
	leaderRepository p2p.LeaderRepository
	dispatcher       *Dispatcher
}

func NewNodeEventHandler(nodeRepo *p2p.NodeRepository, leaderRepo *p2p.LeaderRepository, dispatcher *Dispatcher) *NodeEventHandler {
	return &NodeEventHandler{
		nodeRepository:   nodeRepo,
		leaderRepository: leaderRepo,
		dispatcher:       dispatcher,
	}
}

//todo conn후 peer에 추가하고 peer 추가됬다는 이벤트를 날려줘야 하나 말아야하나 고민이슈
func (neh *NodeEventHandler) HandleConnCreateEvent(event p2p.ConnectionCreatedEvent) {
	id := p2p.NodeId(event.ID)
	address := event.Address
	node := p2p.NewNode(address, id)
	neh.nodeRepository.Save(*node)
	if neh.leaderRepository.GetLeader() == nil {
		err := neh.dispatcher.RequestLeaderInfo(*node)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

//todo conn disconnect event 구현
func (neh *NodeEventHandler) HandleConnDisconnectEvent(event p2p.ConnectionDisconnectedEvent) {
	panic("need to implement")
}

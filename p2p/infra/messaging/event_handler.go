package messaging

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
)

type NodeEventHandler struct {
	nodeRepository   p2p.NodeRepository
	leaderRepository p2p.LeaderRepository
}

func NewNodeEventHandler(nodeRepo *leveldb.NodeRepository, leaderRepo *leveldb.LeaderRepository) *NodeEventHandler {
	return &NodeEventHandler{
		nodeRepository:   nodeRepo,
		leaderRepository: leaderRepo,
	}
}

//todo conn후 peer에 추가하고 peer 추가됬다는 이벤트를 날려줘야 하나 말아야하나 고민이슈
//p2p 에서 노드가 추가된 뒤에 노드가 추가되었다는 사실을 알리는 식으로 doc/ 에 시나리오 작성해 두었습니다.
//참고 바랍니다.
func (n *NodeEventHandler) HandleConnCreatedEvent(event p2p.ConnectionCreatedEvent) {

	if event.ID == "" || event.Address == "" {
		return
	}

	node := p2p.NewNode(event.Address, p2p.NodeId(event.ID))

	n.nodeRepository.Save(*node)
}

//todo conn disconnect event 구현
func (n *NodeEventHandler) HandleConnDisconnectEvent(event p2p.ConnectionDisconnectedEvent) {

	if event.ID == "" {
		return
	}

	nodeId := p2p.NodeId(event.ID)

	n.nodeRepository.Remove(nodeId)
}

//save Leader when LeaderReceivedEvent Detected
func (n *NodeEventHandler) HandleLeaderUpdatedEvent(event p2p.LeaderUpdatedEvent) {

	if event.ID == "" {
		return
	}

	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{Id: event.ID},
	}

	n.leaderRepository.SetLeader(leader)
}

//save all nodes when NodeListReceivedEvent Detected
//node는 각자가 aggregate이기 때문에 aggregate가 동시에 update되는 event는 없습니다.
//event 1개가 aggregate1개를 변화시키는 것

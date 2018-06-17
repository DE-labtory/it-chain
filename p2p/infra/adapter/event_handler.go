package adapter

import (
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
)

type EventHandler struct {
	nodeApi          *NodeApi
}

type NodeApi struct{} //p2p.nodeApi


func NewEventHandler(nodeRepo *WriteOnlyNodeRepository, leaderRepo *WriteOnlyLeaderRepository, nodeApi *NodeApi) *EventHandler {

	return &EventHandler{
		nodeApi:          nodeApi,
	}
}

type WriteOnlyNodeRepository struct{} //p2p.NodeRepository
func (nodeRepository WriteOnlyNodeRepository) Save(data p2p.Node) (err error){
	return err
}
func (nodeRepository WriteOnlyNodeRepository) Remove(nodeId p2p.NodeId) (err error){
	return err
}
type WriteOnlyLeaderRepository struct{} //p2p.LeaderRepository
func (leaderRepository WriteOnlyLeaderRepository) SetLeader(leader p2p.Leader){}

type RepositoryProjector struct {
	nodeRepository   *WriteOnlyNodeRepository
	leaderRepository *WriteOnlyLeaderRepository
}

//todo conn후 peer에 추가하고 peer 추가됬다는 이벤트를 날려줘야 하나 말아야하나 고민이슈
//p2p 에서 노드가 추가된 뒤에 노드가 추가되었다는 사실을 알리는 식으로 doc/ 에 시나리오 작성해 두었습니다.
//참고 바랍니다.
func (projector *RepositoryProjector) HandleConnCreatedEvent(event p2p.ConnectionCreatedEvent) {

	if event.ID == "" || event.Address == "" {
		return
	}

	node := *p2p.NewNode(event.Address, p2p.NodeId{Id: event.ID})
	err := projector.nodeRepository.Save(node)

	if err != nil {
		log.Println(err)
	}
}

//todo conn disconnect event 구현
func (projector *RepositoryProjector) HandleConnDisconnectedEvent(event p2p.ConnectionDisconnectedEvent) {

	if event.ID == "" {
		return
	}

	err := projector.nodeRepository.Remove(p2p.NodeId{Id: event.ID})

	if err != nil {
		log.Println(err)
	}
}

//save Leader when LeaderReceivedEvent Detected
func (projector *RepositoryProjector) HandleLeaderUpdatedEvent(event p2p.LeaderUpdatedEvent) {

	if event.ID == "" {
		return
	}

	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{Id: event.ID},
	}

	projector.leaderRepository.SetLeader(leader)
}

func (projector *RepositoryProjector) HandlerNodeCreatedEvent(event p2p.NodeCreatedEvent) {

	if event.ID == "" || event.IpAddress == "" {
		return
	}

	node := p2p.NewNode(event.IpAddress, p2p.NodeId{Id: event.ID})
	projector.nodeRepository.Save(*node)
}

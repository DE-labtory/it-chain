package adapter

import (
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
	"errors"
)

var ErrEmptyAddress = errors.New("empty address proposed")
type EventHandler struct {
	nodeApi          *NodeApi
}

type NodeApi struct{} //p2p.nodeApi


func NewEventHandler(nodeApi *NodeApi) *EventHandler {
	return &EventHandler{
		nodeApi:          nodeApi,
	}
}

type RepositoryProjector struct {
	NodeRepository   p2p.NodeRepository
	LeaderRepository p2p.LeaderRepository
}

//todo conn후 peer에 추가하고 peer 추가됬다는 이벤트를 날려줘야 하나 말아야하나 고민이슈
//p2p 에서 노드가 추가된 뒤에 노드가 추가되었다는 사실을 알리는 식으로 doc/ 에 시나리오 작성해 두었습니다.
//참고 바랍니다.
func NewRepositoryProjector(nodeRepository p2p.NodeRepository, leaderRepository p2p.LeaderRepository) *RepositoryProjector{
	return &RepositoryProjector{
		NodeRepository:nodeRepository,
		LeaderRepository:leaderRepository,
	}
}
func (projector *RepositoryProjector) HandleConnCreatedEvent(event p2p.ConnectionCreatedEvent) error {

	if event.ID == "" {
		return ErrEmptyNodeId
	}
	if event.Address == "" {
		return ErrEmptyAddress
	}

	node := *p2p.NewNode(event.Address, p2p.NodeId{Id: event.ID})
	err := projector.NodeRepository.Save(node)

	if err != nil {
		log.Println(err)
	}
	return nil
}

//todo conn disconnect event 구현
func (projector *RepositoryProjector) HandleConnDisconnectedEvent(event p2p.ConnectionDisconnectedEvent) error {

	if event.ID == "" {
		return ErrEmptyNodeId
	}

	err := projector.NodeRepository.Remove(p2p.NodeId{Id: event.ID})

	if err != nil {
		log.Println(err)
	}
	return nil
}

//save Leader when LeaderReceivedEvent Detected
func (projector *RepositoryProjector) HandleLeaderUpdatedEvent(event p2p.LeaderUpdatedEvent) error {

	if event.ID == "" {
		return ErrEmptyNodeId
	}

	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{Id: event.ID},
	}

	projector.LeaderRepository.SetLeader(leader)
	return nil
}

func (projector *RepositoryProjector) HandlerNodeCreatedEvent(event p2p.NodeCreatedEvent) error {

	if event.ID == "" {
		return ErrEmptyNodeId
	}

	if event.IpAddress == "" {
		return ErrEmptyAddress
	}

	node := p2p.NewNode(event.IpAddress, p2p.NodeId{Id: event.ID})
	projector.NodeRepository.Save(*node)

	return nil
}

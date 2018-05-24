package peer

// 피어 테이블 구조체는
// 자신이 포함된 chain에서 자신의 peer 정보와 리더 정보를, 그리고 리더를 변경하고 가져오는 메소드를 가진다.
// 즉, 피어 테이블의 목적은 전체 db와 별게로 하나의 피어의 입장에서 보다 주관적으로 로직을 수행하기 위함으로 보인다.

type NodeTable struct {
	leader     *Node
	myInfo     *Node
	repository NodeRepository
}

func NewNodeTableService(nodeRepo NodeRepository, myinfo *Node) *NodeTable {
	nodeRepo.Save(*myinfo)
	return &NodeTable{
		leader:     nil,
		myInfo:     myinfo,
		repository: nodeRepo,
	}
}

func (pts *NodeTable) SetLeader(peer *Node) error {
	// todo err handle
	find, _ := pts.repository.FindById(peer.Id)
	if find == nil {

	}
	pts.leader = find
	return nil
}

func (pts *NodeTable) GetLeader() *Node {
	return pts.leader
}

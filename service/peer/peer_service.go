package peer

//peer 최상위 service
type PeerService interface{

	//peer table 조회
	GetPeerTable() PeerTable

	//peer info 찾기
	GetPeerInfoByPeerID(peerID string) *PeerInfo

	//peer info
	PushPeerTable(peerIDs []string)

	//update peerTable
	UpdatePeerTable(peerTable PeerTable)

	//Add peer
	AddPeerInfo(peerInfo *PeerInfo)

	//Request Peer Info
	RequestPeerInfo(host string, port string) *PeerInfo

	BroadCastPeerTable(interface{})
}
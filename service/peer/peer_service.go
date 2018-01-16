package peer

import "it-chain/service/domain"

//peer 최상위 service
type PeerService interface{

	//peer table 조회
	GetPeerTable() domain.PeerTable

	//peer info 찾기
	GetPeerInfoByPeerID(peerID string) *domain.PeerInfo

	//peer info
	PushPeerTable(peerIDs []string)
}
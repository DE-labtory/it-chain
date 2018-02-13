package service

import "it-chain/domain"

//peer 최상위 service
type PeerService interface{

	//peer table 조회
	GetPeerTable() domain.PeerTable

	//peer info 찾기
	GetPeerByPeerID(peerID string) *domain.Peer

	//peer info
	PushPeerTable(peerIDs []string)

	//update peerTable
	UpdatePeerTable(peerTable domain.PeerTable)

	//Add peer
	AddPeer(Peer *domain.Peer)

	//Request Peer Info
	RequestPeer(host string, port string) *domain.Peer

	BroadCastPeerTable(interface{})

	GetLeader() *domain.Peer
}
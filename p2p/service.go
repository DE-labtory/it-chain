package p2p

type ICommunicationService interface {
	Dial(ipAddress string) error
	DeliverPLTable(connectionId string, peerLeaderTable PLTable) error
}

type IPeerService interface {
	Save(peer Peer) error
	Remove(peerId PeerId) error
}

type IPLTableService interface {
	GetPLTableFromCommand(command GrpcReceiveCommand) (PLTable, error)
}
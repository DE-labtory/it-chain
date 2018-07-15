package blockchain

type SyncService interface {
	SyncWithPeer(peer Peer) error
	SyncedCheck(peer Peer) (IsSynced, error)
	Construct(peer Peer) error
}

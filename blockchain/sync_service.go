package blockchain

type SyncService interface {
	SyncWithPeer(peer Peer) error
	syncedCheck(peer Peer) (IsSynced, error)
	construct(peer Peer) error
}

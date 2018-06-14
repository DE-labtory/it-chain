package messaging

import "github.com/it-chain/it-chain-Engine/p2p/api"

type CommandHandler struct {
	leaderApi api.LeaderApi
	nodeApi   api.NodeApi
}

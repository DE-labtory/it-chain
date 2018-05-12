package peer

import "github.com/it-chain/it-chain-Engine/conf"

func init() {
	bootNodeIp := conf.GetConfiguration().Common.BootNodeIp
	myIp := conf.GetConfiguration().Common.NodeIp

	if bootNodeIp == myIp {
		err := leaderSelectionApi.messageProducer.LeaderUpdateEvent(*myInfo)
		if err != nil {
			//todo mq error handle 해야하나?
		}
	}
}

package blockchain

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCreateNewBlockChainTest(t *testing.T){

	var blockChains = CreateNewBlockChain(defaultChannelName,defaultPeerId)
	assert.Equal(t,0,len(blockChains.Blocks))
	assert.Equal(t,defaultPeerId,blockChains.Header.PeerID)
	assert.Equal(t,defaultChannelName,blockChains.Header.ChannelName)
}

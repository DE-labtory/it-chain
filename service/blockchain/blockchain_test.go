package blockchain

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"fmt"
)

func TestCreateNewBlockChainTest(t *testing.T){

	var blockChains = CreateNewBlockChain(defaultChannelName,defaultPeerId)
	assert.Equal(t,0,len(blockChains.blocks))
	assert.Equal(t,defaultPeerId,blockChains.header.peerID)
	assert.Equal(t,defaultChannelName,blockChains.header.channelName)
}

package adapter_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/p2p/infra/adapter"
	"github.com/magiconair/properties/assert"
)

type MockNodeRepository struct {}
func (mockNodeRepository MockNodeRepository) Save(){

}
func (mockNodeRepository MockNodeRepository) Remove(){

}
func (mockNodeRepository MockNodeRepository) FindByInd(){

}
func (mockNodeRepository MockNodeRepository) FindAll(){

}
type MockLeaderRepository struct {

}
type MockNodeApi struct {

}
func TestNewEventHandler(t *testing.T) {
	nodeRepository := MockNodeRepository{}
	leaderRepository := MockLeaderRepository{}
	nodeApi := MockNodeApi{}

	eventHandler := adapter.NewEventHandler(nodeRepository, leaderRepository, nodeApi)

	assert.Equal(t, eventHandler.nodeRepository, )

}

func TestEventHandler_HandleConnCreatedEvent(t *testing.T) {

}
func TestEventHandler_HandleConnDisconnectedEvent(t *testing.T) {

}
func TestEventHandler_HandleLeaderUpdatedEvent(t *testing.T) {

}
func TestEventHandler_HandlerNodeCreatedEvent(t *testing.T) {

}

func SetUpEventHandler(){

}
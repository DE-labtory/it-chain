package rabbitmq

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/msg"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestMessageApi_BroadCastMsg(t *testing.T) {

	//given
	messaging := rabbitmq.Connect("amqp://guest:guest@localhost:5672/")
	defer messaging.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	err := messaging.Consume(topic.ConsensusMessagePublishEvent.String(), func(delivery amqp.Delivery) {
		fmt.Println("waiting")
		ReceivedMsg := &event.ConsensusMessagePublishEvent{}
		json.Unmarshal(delivery.Body, ReceivedMsg)

		//then
		assert.Equal(t, []string{"1", "2"}, ReceivedMsg.Ids)
		wg.Done()

	})

	assert.NoError(t, err)

	mApi := NewRabbitmqMessageService(messaging.Publish)

	message := msg.PreprepareMsg{}
	representatives := []*consensus.Representative{&consensus.Representative{"1"}, &consensus.Representative{"2"}}

	//when
	mApi.BroadCastMsg(message, representatives)

	wg.Wait()
}

func TestMessageApi_ConfirmedBlock(t *testing.T) {

	//given
	messaging := rabbitmq.Connect("amqp://guest:guest@localhost:5672/")
	defer messaging.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	blockData := []byte("my block")

	err := messaging.Consume(topic.BlockConfirmEvent.String(), func(delivery amqp.Delivery) {
		fmt.Println("waiting")
		blockConfirmEvent := &event.BlockConfirmEvent{}
		json.Unmarshal(delivery.Body, blockConfirmEvent)

		//then
		assert.Equal(t, blockData, blockConfirmEvent.Block)
		wg.Done()
	})

	assert.NoError(t, err)

	mApi := NewRabbitmqMessageService(messaging.Publish)

	//when
	mApi.ConfirmedBlock(blockData)
}

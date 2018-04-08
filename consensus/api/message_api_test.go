package api

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/msg"
	m "github.com/it-chain/it-chain-Engine/messaging"
	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/it-chain/it-chain-Engine/messaging/topic"
	"github.com/stretchr/testify/assert"
)

func TestMessageApi_BroadCastMsg(t *testing.T) {

	//given
	messaging := m.NewMessaging("amqp://guest:guest@localhost:5672/")
	messaging.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)

	msgs, _ := messaging.Consume(topic.ConsensusMessagePublishEvent.String())

	go func() {
		fmt.Println("waiting")
		for data := range msgs {
			ReceivedMsg := &event.ConsensusMessagePublishEvent{}
			json.Unmarshal(data.Body, ReceivedMsg)

			//then
			assert.Equal(t, []string{"1", "2"}, ReceivedMsg.Ids)
			wg.Done()
		}
	}()

	mApi := NewMessageApi(messaging.Publish)

	message := msg.PreprepareMsg{}
	representatives := []*consensus.Representative{&consensus.Representative{"1"}, &consensus.Representative{"2"}}

	//when
	mApi.BroadCastMsg(message, representatives)

	wg.Wait()
}

func TestMessageApi_ConfirmedBlock(t *testing.T) {

	//given
	messaging := m.NewMessaging("amqp://guest:guest@localhost:5672/")
	messaging.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)

	blockData := []byte("my block")

	msgs, _ := messaging.Consume(topic.BlockConfirmEvent.String())

	go func() {
		fmt.Println("waiting")
		for data := range msgs {
			blockConfirmEvent := &event.BlockConfirmEvent{}
			json.Unmarshal(data.Body, blockConfirmEvent)

			//then
			assert.Equal(t, blockData, blockConfirmEvent.Block)
			wg.Done()
		}
	}()

	mApi := NewMessageApi(messaging.Publish)

	//when
	mApi.ConfirmedBlock(blockData)
}

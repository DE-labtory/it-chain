package rabbitmq

import (
	"fmt"
	"sync"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestMessaging_Start(t *testing.T) {

	message := Connect("amqp://guest:guest@localhost:5672/")
	defer message.Close()
}

func TestMessaging_Publish(t *testing.T) {

	message := Connect("amqp://guest:guest@localhost:5672/")
	defer message.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	err := message.Consume("asd", func(delivery amqp.Delivery) {
		fmt.Println("received data", delivery)
		wg.Done()
	})

	assert.NoError(t, err)

	fmt.Println("waiting1")
	err = message.Publish("asd", []byte("zxc"))

	if err != nil {
		assert.NoError(t, err)
	}

	wg.Wait()
}

func TestMessaging_MultiPublishAndConsume(t *testing.T) {

	message := Connect("amqp://guest:guest@localhost:5672/")
	defer message.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	err := message.Consume("asd", func(delivery amqp.Delivery) {
		assert.Equal(t, delivery.Body, []byte("zxc"))
		wg.Done()
	})

	assert.NoError(t, err)

	err = message.Consume("asd1", func(delivery amqp.Delivery) {
		assert.Equal(t, delivery.Body, []byte("zxc"))
		wg.Done()
	})

	assert.NoError(t, err)

	err = message.Publish("asd", []byte("zxc"))

	if err != nil {
		assert.NoError(t, err)
	}

	err = message.Publish("asd1", []byte("zxc"))

	if err != nil {
		assert.NoError(t, err)
	}

	wg.Wait()
}

func TestMultiConsumer(t *testing.T) {

	message := Connect("amqp://guest:guest@localhost:5672/")
	defer message.Close()

	wg := sync.WaitGroup{}

	wg.Add(1)
	err := message.Consume("asd", func(delivery amqp.Delivery) {
		fmt.Println("1")
		wg.Done()
	})

	assert.NoError(t, err)

	message2 := Connect("amqp://guest:guest@localhost:5672/")
	defer message2.Close()

	wg.Add(1)
	err = message2.Consume("asd", func(delivery amqp.Delivery) {
		fmt.Println("2")
		wg.Done()
	})

	assert.NoError(t, err)

	message3 := Connect("amqp://guest:guest@localhost:5672/")
	defer message3.Close()

	err = message3.Publish("asd", []byte("asd"))

	assert.NoError(t, err)

	wg.Wait()

}

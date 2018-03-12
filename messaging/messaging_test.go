package messaging

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"sync"
	"fmt"
)

func TestMessaging_Start(t *testing.T) {

	message := NewMessaging("amqp://guest:guest@localhost:5672/")
	assert.NotPanics(t,message.Start)
}

func TestMessaging_Start2(t *testing.T) {

	message := NewMessaging("amqp://guest:guest@localhost:5672/")
	assert.NotPanics(t,message.Start)
	err := message.Publish("asd",[]byte("zxc"))

	if err != nil{
		assert.NoError(t,err)
	}
}

func TestMessaging_Publish(t *testing.T) {

	message := NewMessaging("amqp://guest:guest@localhost:5672/")
	message.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)

	msg, err := message.Consume("asd")

	if err != nil{

	}

	go func (){
		fmt.Println("waiting")
		for data := range msg{
			fmt.Println("received data", data)
			wg.Done()
		}
	}()

	fmt.Println("waiting1")
	err = message.Publish("asd",[]byte("zxc"))

	if err != nil{
		assert.NoError(t,err)
	}

	wg.Wait()
}

func TestMessaging_MultiPublishAndConsume(t *testing.T) {

	message := NewMessaging("amqp://guest:guest@localhost:5672/")
	message.Start()

	wg := sync.WaitGroup{}
	wg.Add(2)

	asdMsg, err := message.Consume("asd")

	if err != nil{

	}

	asd1Msg, err := message.Consume("asd1")

	if err != nil{

	}

	go func (){
		for data := range asdMsg{
			assert.Equal(t,data.Body,[]byte("zxc"))
			wg.Done()
		}
	}()

	go func (){
		for data := range asd1Msg{
			assert.Equal(t,data.Body,[]byte("zxc"))
			wg.Done()
		}
	}()

	err = message.Publish("asd",[]byte("zxc"))

	if err != nil{
		assert.NoError(t,err)
	}

	err = message.Publish("asd1",[]byte("zxc"))

	if err != nil{
		assert.NoError(t,err)
	}

	wg.Wait()
}
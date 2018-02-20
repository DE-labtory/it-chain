package publisher

import (
	"it-chain/common"
	"it-chain/network/comm"
	"errors"
	"it-chain/auth"
	"sync"
)

var logger_event_publisher = common.GetLogger("message_publisher.go")

//message를 받으면 관심이 있는 subscriber에게 전달하는 역활을 한다.
type MessagePublisher struct {
	subscribers map[string]func(message comm.OutterMessage)
	crpyto      auth.Crypto
	lock        *sync.Mutex
}

//todo Message Publisher를 connectionManager가 관리하게 할지 고민
//topic의 일치성을 위해 처음에 topic의 list를 받고 이 list에 없는 topic은 등록불가하다.
func NewMessagePublisher(crpyto auth.Crypto) *MessagePublisher{

	return &MessagePublisher{
		subscribers: make(map[string]func(message comm.OutterMessage)),
		crpyto: crpyto,
		lock: &sync.Mutex{},
	}
}

//subscriber를 등록한다.
func (mp *MessagePublisher) AddSubscriber(name string, subfunc func(message comm.OutterMessage)) error{

	mp.lock.Lock()
	defer mp.lock.Unlock()

	_, ok := mp.subscribers[name]

	if ok{
		return errors.New("already subscribed function")
	}

	mp.subscribers[name] = subfunc

	return nil
}

func (mp *MessagePublisher) ReceivedMessageHandle(message comm.OutterMessage){

	defer recover()

	envelope := message.Envelope

	if envelope == nil{
		logger_event_publisher.Info("message is nil", message)
		return
	}

	_, pub, err := mp.crpyto.GetKey()

	if err != nil{
		logger_event_publisher.Infof("failed to load key: %s", err.Error())
		return
	}

	vaild, err := mp.crpyto.Verify(pub,envelope.Signature,envelope.Payload,nil)

	if !vaild || err !=nil{
		logger_event_publisher.Info("failed to verify message")
		return
	}

	msg, err := envelope.GetMessage()

	message.Message = msg

	if err != nil{
		logger_event_publisher.Info("failed to Unmarshal message:", message)
		return
	}

	mp.lock.Lock()
	defer mp.lock.Unlock()

	for _, subFunc := range mp.subscribers{
		subFunc(message)
	}
}


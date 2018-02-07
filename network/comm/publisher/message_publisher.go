package publisher

import (
	"github.com/asaskevich/EventBus"
	"it-chain/common"
	"it-chain/network/comm"
	pb "it-chain/network/protos"
	"github.com/golang/protobuf/proto"
	"errors"
	"it-chain/auth"
)

var logger_event_publisher = common.GetLogger("message_publisher.go")

//message를 받으면 관심이 있는 subscriber에게 전달하는 역활을 한다.
type MessagePublisher struct{
	bus 		EventBus.Bus
	topicMap    map[string]string
	crpyto      auth.Crypto
}

//todo Message Publisher를 connectionManager가 관리하게 할지 고민
//todo signer를 받아야함
//topic의 일치성을 위해 처음에 topic의 list를 받고 이 list에 없는 topic은 등록불가하다.
func NewMessagePublisher(messageTypes []string, crpyto auth.Crypto) *MessagePublisher{

	topicMap := make(map[string]string)

	for _,messageType := range messageTypes{
		topicMap[messageType] = messageType
	}

	return &MessagePublisher{
		bus: EventBus.New(),
		topicMap: topicMap,
		crpyto: crpyto,
	}
}

//subscriber를 등록한다.
func (mp *MessagePublisher) AddSubscriber(topic string, subfunc func(message comm.OutterMessage), transactional bool) error{

	t, ok := mp.topicMap[topic]

	if ok {
		err := mp.bus.SubscribeAsync(t, subfunc, transactional)

		if err != nil {
			logger_event_publisher.Error("failed to add subscriber", err.Error())
			return err
		}

		logger_event_publisher.Infoln("new message subscriber added")
		return nil
	}

	logger_event_publisher.Error("failed to add subscriber: invalid topic")
	return errors.New("invaild topic")
}

//todo invaild message검증 및 전파
//todo panic으로 부터 recover
//todo topic을 미리 만들어 놓는 방식에 대해서 고민 해야함
func (mp *MessagePublisher) ReceivedMessageHandle(message comm.OutterMessage){

	envelope := message.Envelope

	if envelope == nil{
		logger_event_publisher.Info("message is nil", message)
		return
	}

	_, pub, err := mp.crpyto.LoadKey()

	if err != nil{
		logger_event_publisher.Infof("failed to load key: %s", err.Error())
		return
	}

	vaild, err := mp.crpyto.Verify(pub,envelope.Signature,envelope.Payload,nil)

	if !vaild || err !=nil{
		logger_event_publisher.Info("failed to verify message")
		return
	}

	m := &pb.Message{}

	err = proto.Unmarshal(envelope.Payload,m)

	if err != nil{
		logger_event_publisher.Info("failed to Unmarshal message:", message)
		return
	}

	message.Message = m
	messageType := m.GetMessageType()
	mt, ok := mp.topicMap[messageType]

	if ok{
		mp.bus.Publish(mt,message)
		logger_event_publisher.Info("message published messageType:", mt)
	}
}


package adapter

// ToDo: 삭제 제안 - junk_sound

//// ToDo: 구현.(gitId:junk-sound)
//type Publish func(exchange string, topic string, data interface{}) (err error)
//
//type GrpcCommandService struct {
//	publish Publish // midgard.client.Publish
//}
//
//
//
//func NewGrpcCommandService(publish Publish) *GrpcCommandService {
//	return &GrpcCommandService{
//		publish: publish,
//	}
//}
//
//func createGrpcDeliverCommand(protocol string, body interface{}) (blockchain.GrpcDeliverCommand, error) {
//
//	data, err := common.Serialize(body)
//	if err != nil {
//		return blockchain.GrpcDeliverCommand{}, err
//	}
//
//	return blockchain.GrpcDeliverCommand{
//		CommandModel: midgard.CommandModel{
//			ID: xid.New().String(),
//		},
//		Recipients: make([]string, 0),
//		Body:       data,
//		Protocol:   protocol,
//	}, err
//}

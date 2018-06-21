package p2p

func init() {
	////myIp := conf.GetConfiguration().Common.NodeIp
	//
	//config := conf.GetConfiguration()
	////create rabbitmq client
	//rabbitmqClient := rabbitmq.Connect(config.Common.Messaging.Url)
	//// todo change node repo and leader repo after managing peerTable
	//nodeRepository := leveldb.NewNodeRepository("path1")
	//leaderRepository := leveldb.NewLeaderRepository("path2")
	//publisher := rabbitmq.Connect("")
	//
	//messageDispatcher := adapter.NewMessageDispatcher(publisher)
	//
	//leaderApi := api.NewLeaderApi(eventRepository, messageDispatcher, myInfo)
	//
	////create amqp Handler
	//eventHandler := adapter.NewNodeEventHandler(nodeRepository, leaderRepository)
	//grpcMessageHandler := adapter.NewGrpcMessageHandler(leaderApi, nodeApi, messageDispatcher)
	//
	//// Subscribe amqp server
	//err1 := rabbitmqClient.Subscribe("Command", "connection.*", eventHandler)
	//
	//if err1 != nil {
	//	panic(err1)
	//}
	//
	//err2 := rabbitmqClient.Subscribe("Command", "message.*", grpcMessageHandler)
	//
	//if err2 != nil {
	//	panic(err2)
	//}
}

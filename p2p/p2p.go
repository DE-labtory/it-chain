/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package p2p

//todo implement it!
func init() {
	//publisher := rabbitmq.Connect("")
	//eventRepository := midgard.NewRepo(store, publisher)
	//publisher :=
	//grpcCommandService := adapter.NewGrpcCommandService(publisher)
	//peerApi := api.NewPeerApi(eventRepository, grpcCommandService)
	////myIp := conf.GetConfiguration().Common.NodeIp
	//
	//config := conf.GetConfiguration()
	////create rabbitmq client
	//rabbitmqClient := rabbitmq.Connect(config.Common.Messaging.Url)
	//// todo change node repo and leader repo after managing peerTable
	//
	//messageDispatcher := adapter.NewMessageDispatcher(publisher)
	//
	//
	////create rabbitmq Handler
	//eventHandler := adapter.NewNodeEventHandler(nodeRepository, leaderRepository)
	//grpcMessageHandler := adapter.NewGrpcMessageHandler(leaderApi, nodeApi, messageDispatcher)
	//
	//// Subscribe rabbitmq server
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

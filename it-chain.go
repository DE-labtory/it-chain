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

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/blockchain"
	blockchainApi "github.com/it-chain/engine/blockchain/api"
	blockchainAdapter "github.com/it-chain/engine/blockchain/infra/adapter"
	blockchainMem "github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/cmd/ivm"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	icodeApi "github.com/it-chain/engine/ivm/api"
	icodeAdapter "github.com/it-chain/engine/ivm/infra/adapter"
	icodeInfra "github.com/it-chain/engine/ivm/infra/git"
	"github.com/it-chain/engine/ivm/infra/tesseract"
	txpoolApi "github.com/it-chain/engine/txpool/api"
	txpoolAdapter "github.com/it-chain/engine/txpool/infra/adapter"
	txpoolBatch "github.com/it-chain/engine/txpool/infra/batch"
	txpoolMem "github.com/it-chain/engine/txpool/infra/mem"
	"github.com/urfave/cli"
	"github.com/it-chain/engine/cmd/connection"
)

const apidbPath = "./api-db"
const dbPath = "./db"

func PrintLogo() {
	fmt.Println(`
	___  _________               ________  ___  ___  ________  ___  ________
	|\  \|\___   ___\            |\   ____\|\  \|\  \|\   __  \|\  \|\   ___  \
	\ \  \|___ \  \_|____________\ \  \___|\ \  \\\  \ \  \|\  \ \  \ \  \\ \  \
	 \ \  \   \ \  \|\____________\ \  \    \ \   __  \ \   __  \ \  \ \  \\ \  \
	  \ \  \   \ \  \|____________|\ \  \____\ \  \ \  \ \  \ \  \ \  \ \  \\ \  \
           \ \__\   \ \__\              \ \_______\ \__\ \__\ \__\ \__\ \__\ \__\\ \__\
	    \|__|    \|__|               \|_______|\|__|\|__|\|__|\|__|\|__|\|__| \|__|
	`)
}

func main() {

	app := cli.NewApp()
	app.Name = "it-chain"
	app.Version = "0.1.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "it-chain",
			Email: "it-chain@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "config",
			Usage: "name for config",
		},
	}
	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, ivm.IcodeCmd())
	app.Commands = append(app.Commands, connection.Cmd())
	app.Action = func(c *cli.Context) error {
		PrintLogo()
		configName := c.String("config")
		conf.SetConfigName(configName)
		return run()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	clearFolder()
	errs := make(chan error, 2)

	rpcServer, rpcClient, configuration, tearDown := initCommon()
	defer tearDown()

	logger.EnableFileLogger(true, configuration.Engine.LogPath)

	defer initApiGateway(configuration, errs)()
	defer initTxPool(configuration, rpcServer, rpcClient)()
	defer initICode(configuration, rpcServer)()
	defer initBlockchain(configuration, rpcServer)()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Println("terminated", <-errs)

	return nil
}

func clearFolder() {

	if err := os.RemoveAll(apidbPath); err != nil {
		logger.Panic(&logger.Fields{"err_msg": err.Error()}, "error while clear folder")
		panic(err)
	}

	if err := os.RemoveAll(dbPath); err != nil {
		logger.Panic(&logger.Fields{"err_msg": err.Error()}, "error while clear folder")
		panic(err)
	}
}

func initCommon() (rpc.Server, rpc.Client, *conf.Configuration, func()) {
	config := conf.GetConfiguration()
	server := rpc.NewServer(config.Engine.Amqp)
	client := rpc.NewClient(config.Engine.Amqp)
	return server, client, config, func() {
		server.Close()
		client.Close()
	}
}

func initApiGateway(config *conf.Configuration, errs chan error) func() {

	ipAddress := config.ApiGateway.Address + ":" + config.ApiGateway.Port

	//set log
	var kitLogger kitlog.Logger
	kitLogger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	kitLogger = kitlog.With(kitLogger, "ts", kitlog.DefaultTimestampUTC)

	// set blockchain
	blockchainDB := apidbPath + "/block"
	CommittedBlockRepo, err := api_gateway.NewBlockRepositoryImpl(blockchainDB)
	if err != nil {
		logger.Panic(&logger.Fields{"err_msg": err.Error()}, "error while init gateway")
		panic(err)
	}

	blockQueryApi := api_gateway.NewBlockQueryApi(CommittedBlockRepo)
	blockEventListener := api_gateway.NewBlockEventListener(CommittedBlockRepo)

	// set ivm
	icodeDB := apidbPath + "/ivm"
	icodeRepo := api_gateway.NewLevelDbMetaRepository(icodeDB)
	icodeQueryApi := api_gateway.NewICodeQueryApi(&icodeRepo)
	icodeEventListener := api_gateway.NewIcodeEventHandler(&icodeRepo)

	//set mux
	mux := http.NewServeMux()
	httpLogger := kitlog.With(kitLogger, "component", "http")

	subscriber := pubsub.NewTopicSubscriber(config.Engine.Amqp, "Event")
	if err := subscriber.SubscribeTopic("block.*", &blockEventListener); err != nil {
		panic(err)
	}
	if err := subscriber.SubscribeTopic("icode.*", &icodeEventListener); err != nil {
		panic(err)
	}

	var handler http.Handler
	{
		handler = api_gateway.NewApiHandler(blockQueryApi, icodeQueryApi, httpLogger)
	}

	http.Handle("/", mux)

	go func() {
		logger.Infof(nil, "[Main] Api-gateway is staring on port:%s", config.ApiGateway.Port)
		errs <- http.ListenAndServe(ipAddress, handler)
	}()

	return func() {
		CommittedBlockRepo.Close()
		os.RemoveAll(apidbPath)
	}
}

func initICode(config *conf.Configuration, server rpc.Server) func() {

	logger.Infof(nil, "[Main] Ivm is staring")

	// git generate
	storeApi := icodeInfra.NewRepositoryService()
	containerService := tesseract.NewContainerService()
	eventService := common.NewEventService(config.Engine.Amqp, "Event")
	api := icodeApi.NewICodeApi(containerService, storeApi, eventService)

	// handler generate
	deployHandler := icodeAdapter.NewDeployCommandHandler(api)
	unDeployHandler := icodeAdapter.NewUnDeployCommandHandler(api)
	icodeExecuteHandler := icodeAdapter.NewIcodeExecuteCommandHandler(api)
	listHandler := icodeAdapter.NewListCommandHandler(api)
	blockCommittedEventHandler := icodeAdapter.NewBlockCommittedEventHandler(api)

	server.Register("ivm.execute", icodeExecuteHandler.HandleTransactionExecuteCommandHandler)
	server.Register("ivm.deploy", deployHandler.HandleDeployCommand)
	server.Register("ivm.undeploy", unDeployHandler.HandleUnDeployCommand)
	server.Register("ivm.list", listHandler.HandleListCommand)

	subscriber := pubsub.NewTopicSubscriber(config.Engine.Amqp, "Event")
	if err := subscriber.SubscribeTopic("block.*", blockCommittedEventHandler); err != nil {
		panic(err)
	}

	return func() {
		iCodeInfos := containerService.GetRunningICodeList()
		for _, iCodeInfo := range iCodeInfos {
			containerService.StopContainer(iCodeInfo.ID)
		}
	}
}

func initTxPool(config *conf.Configuration, server rpc.Server, client rpc.Client) func() {

	logger.Infof(nil, "[Main] Txpool is staring")

	//todo get id from pubkey
	tmpPeerID := "tmp peer 1"
	transactionRepo := txpoolMem.NewTransactionRepository()
	blockProposalService := txpoolAdapter.NewBlockProposalService(client, transactionRepo, config.Engine.Mode)
	txApi := txpoolApi.NewTransactionApi(tmpPeerID, transactionRepo)
	txCommandHandler := txpoolAdapter.NewTxCommandHandler(txApi)
	txpoolBatch.GetTimeOutBatcherInstance().Run(blockProposalService.ProposeBlock, (time.Duration(config.Txpool.TimeoutMs) * time.Millisecond))

	err := server.Register("transaction.create", txCommandHandler.HandleTxCreateCommand)

	if err != nil {
		panic(err)
	}

	return func() {}
}

func initBlockchain(config *conf.Configuration, server rpc.Server) func() {

	logger.Infof(nil, "[Main] Blockchain is staring")

	publisherId := "publisher.1"
	blockRepo, err := blockchainMem.NewBlockRepository(dbPath)

	if err != nil {
		panic(err)
	}

	eventService := common.NewEventService(config.Engine.Amqp, "Event")
	blockPool := blockchain.NewBlockPool()

	blockApi, err := blockchainApi.NewBlockApi(publisherId, blockRepo, eventService, blockPool)
	if err != nil {
		panic(err)
	}

	err = blockApi.CommitGenesisBlock(config.Blockchain.GenesisConfPath)
	if err != nil {
		panic(err)
	}

	// TODO: Change with real consensus service
	consensusService := mock.ConsensusService{}

	blockProposeHandler := blockchainAdapter.NewBlockProposeCommandHandler(blockApi, consensusService, config.Engine.Mode)
	server.Register("block.propose", blockProposeHandler.HandleProposeBlockCommand)

	return func() {
		os.RemoveAll(dbPath)
	}
}

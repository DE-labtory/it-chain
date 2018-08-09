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
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/it-chain/engine/api_gateway"
	blockchainApi "github.com/it-chain/engine/blockchain/api"
	blockchainAdapter "github.com/it-chain/engine/blockchain/infra/adapter"
	blockchainMem "github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/it-chain/engine/cmd/icode"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	icodeApi "github.com/it-chain/engine/icode/api"
	icodeAdapter "github.com/it-chain/engine/icode/infra/adapter"
	icodeInfra "github.com/it-chain/engine/icode/infra/git"
	"github.com/it-chain/engine/icode/infra/tesseract"
	txpoolApi "github.com/it-chain/engine/txpool/api"
	txpoolAdapter "github.com/it-chain/engine/txpool/infra/adapter"
	txpoolBatch "github.com/it-chain/engine/txpool/infra/batch"
	txpoolMem "github.com/it-chain/engine/txpool/infra/mem"
	"github.com/urfave/cli"
)

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
	app.Commands = append(app.Commands, icode.IcodeCmd())
	app.Action = func(c *cli.Context) error {
		PrintLogo()
		configName := c.String("config")
		conf.SetConfigName(configName)
		return start()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func start() error {

	configuration := conf.GetConfiguration()
	logger.EnableFileLogger(true, configuration.Engine.LogPath)

	ip4 := configuration.GrpcGateway.Address + ":" + configuration.GrpcGateway.Port
	ln, err := net.Listen("tcp", ip4)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't listen on %q: %s\n", ip4, err)
		return err
	}

	err = ln.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't stop listening on %q: %s\n", ip4, err)
		return err
	}

	errs := make(chan error, 2)

	server, client, config := initCommon()
	gatewayTearDown := initGateway(errs)
	txPoolTearDown := initTxPool(server, client, config)
	iCodeTearDown := initIcode(server, config)
	peerTearDown := initPeer()
	blockChainTearDown := initBlockchain(server, config)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	select {
	case <-errs:
		e := gatewayTearDown()
		if e != nil {
			logger.Error(nil, "error while tear down gateway")
		}
		e = txPoolTearDown()
		if e != nil {
			logger.Error(nil, "error while tear down txpool")
		}
		e = iCodeTearDown()
		if e != nil {
			logger.Error(nil, "error while tear down icode")
		}
		e = peerTearDown()
		if e != nil {
			logger.Error(nil, "error while tear down peer")
		}
		e = blockChainTearDown()
		if e != nil {
			logger.Error(nil, "error while tear down block chain")
		}
	default:

	}

	return nil
}

//todo other way to inject each query Api to component
var blockQueryApi api_gateway.BlockQueryApi

func initCommon() (rpc.Server, rpc.Client, *conf.Configuration) {
	config := conf.GetConfiguration()
	server := rpc.NewServer(config.Engine.Amqp)
	client := rpc.NewClient(config.Engine.Amqp)
	return server, client, config
}

func initGateway(errs chan error) func() error {

	log.Println("gateway is running...")

	config := conf.GetConfiguration()
	ipAddress := config.ApiGateway.Address + ":" + config.ApiGateway.Port

	//set log
	var kitLogger kitlog.Logger
	kitLogger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	kitLogger = kitlog.With(kitLogger, "ts", kitlog.DefaultTimestampUTC)

	//set subscriber
	subscriber := pubsub.NewTopicSubscriber(config.Engine.Amqp, "Event")

	//set blockchain tesseract and repo

	blockchainDB := "./.test/blockchain"
	BlockPoolRepo := api_gateway.NewBlockPoolRepository()
	CommittedBlockRepo, err := api_gateway.NewCommitedBlockRepositoryImpl(blockchainDB)
	blockQueryApi = api_gateway.NewBlockQueryApi(BlockPoolRepo, CommittedBlockRepo)
	blockEventListener := api_gateway.NewBlockEventListener(BlockPoolRepo, CommittedBlockRepo)

	//set icode tesseract and repo

	icodeDb := "./.test/icode"
	icodeRepo := api_gateway.NewLevelDbMetaRepository(icodeDb)
	metaQueryApi := api_gateway.NewICodeQueryApi(&icodeRepo)
	metaEventListener := api_gateway.NewIcodeEventHandler(&icodeRepo)

	if err != nil {
		logger.Panic(&logger.Fields{"err_msg": err.Error()}, "error while init gateway")
		panic(err)
	}

	//set mux
	mux := http.NewServeMux()
	httpLogger := kitlog.With(kitLogger, "component", "http")

	err = subscriber.SubscribeTopic("block.*", &blockEventListener)
	err = subscriber.SubscribeTopic("meta.*", &metaEventListener)

	if err != nil {
		panic(err)
	}

	mux.Handle("/blocks", api_gateway.BlockchainApiHandler(blockQueryApi, httpLogger))
	mux.Handle("/metas", api_gateway.ICodeApiHandler(metaQueryApi, httpLogger))
	http.Handle("/", mux)

	go func() {
		log.Println("transport", "http", "address", ipAddress, "msg", "listening")
		errs <- http.ListenAndServe(ipAddress, nil)
	}()
	tearDown := func() error {
		CommittedBlockRepo.Close()
		e := os.RemoveAll("./.test")
		return e
	}
	return tearDown
}

func initIcode(server rpc.Server, config *conf.Configuration) func() error {

	log.Println("icode is running...")
	//publisher := pubsub.NewTopicPublisher(config.Engine.Amqp, "Command")

	// tesseract generate
	//commandService := icodeAdapter.NewCommandService(publisher.Publish)

	// git generate
	storeApi := icodeInfra.NewRepositoryService()
	containerService := tesseract.NewContainerService()
	eventService := common.NewEventService(config.Engine.Amqp, "Event")
	api := icodeApi.NewICodeApi(containerService, storeApi, eventService)

	// handler generate
	deployHandler := icodeAdapter.NewDeployCommandHandler(api)
	unDeployHandler := icodeAdapter.NewUnDeployCommandHandler(api)
	icodeExecuteHandler := icodeAdapter.NewIcodeExecuteCommandHandler(api)

	server.Register("icode.execute", icodeExecuteHandler.HandleTransactionExecuteCommandHandler)
	server.Register("icode.deploy", deployHandler.HandleDeployCommand)
	server.Register("icode.undeploy", unDeployHandler.HandleUnDeployCommand)

	tearDown := func() error {
		server.Close()
		return nil
	}
	return tearDown
}

func initPeer() func() error {
	return nil
}

func initTxPool(server rpc.Server, client rpc.Client, config *conf.Configuration) func() error {

	log.Println("txpool is running...")

	//todo get id from pubkey
	tmpPeerID := "tmp peer 1"

	transactionRepo := txpoolMem.NewTransactionRepository()

	//tesseract
	blockProposalService := txpoolAdapter.NewBlockProposalService(client, transactionRepo, config.Engine.Mode)

	//infra
	txApi := txpoolApi.NewTransactionApi(tmpPeerID, transactionRepo)
	txCommandHandler := txpoolAdapter.NewTxCommandHandler(txApi)

	//10초마다 block propose
	txpoolBatch.GetTimeOutBatcherInstance().Run(blockProposalService.ProposeBlock, (time.Duration(config.Txpool.TimeoutMs) * time.Millisecond))

	err := server.Register("transaction.create", txCommandHandler.HandleTxCreateCommand)

	if err != nil {
		panic(err)
	}

	tearDown := func() error {
		server.Close()
		client.Close()
		return nil
	}
	return tearDown
}

func initConsensus() func() error {
	return nil
}

func initBlockchain(server rpc.Server, config *conf.Configuration) func() error {

	log.Println("blockchain is running...")

	publisherId := "publisher.1"

	blockRepo, err := blockchainMem.NewBlockRepository("./blockchain/db")

	if err != nil {
		panic(err)
	}

	eventService := common.NewEventService(config.Engine.Amqp, "Event")

	// api
	blockApi, err := blockchainApi.NewBlockApi(publisherId, blockRepo, eventService)
	if err != nil {
		panic(err)
	}

	// infra
	blockProposeHandler := blockchainAdapter.NewBlockProposeCommandHandler(blockApi, config.Engine.Mode)

	server.Register("block.propose", blockProposeHandler.HandleProposeBlockCommand)

	tearDown := func() error {
		server.Close()
		blockRepo.Close()
		e := os.RemoveAll("./blockchain/db")
		return e
	}
	return tearDown
}

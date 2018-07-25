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
	"github.com/it-chain/engine/cmd/icode"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/core/eventstore"
	icodeApi "github.com/it-chain/engine/icode/api"
	icodeAdapter "github.com/it-chain/engine/icode/infra/adapter"
	icodeInfra "github.com/it-chain/engine/icode/infra/git"
	icodeService "github.com/it-chain/engine/icode/infra/service"
	"github.com/it-chain/engine/txpool"
	txpoolApi "github.com/it-chain/engine/txpool/api"
	txpoolAdapter "github.com/it-chain/engine/txpool/infra/adapter"
	txpoolBatch "github.com/it-chain/engine/txpool/infra/batch"
	"github.com/it-chain/tesseract"
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
		eventstore.InitDefault()
		return start()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func start() error {

	configuration := conf.GetConfiguration()
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

	initGateway(errs)
	initTxPool()
	initIcode()
	initPeer()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Println("terminated", <-errs)

	return nil
}

//todo other way to inject each query Api to component
var txQueryApi api_gateway.TransactionQueryApi

func initGateway(errs chan error) error {

	log.Println("gateway is running...")

	config := conf.GetConfiguration()
	ipAddress := config.ApiGateway.Address + ":" + config.ApiGateway.Port

	//set log
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	//set service and repo
	dbPath := "./.test"
	subscriber := pubsub.NewTopicSubscriber(config.Engine.Amqp, "Event")

	repo := api_gateway.NewTransactionRepository(dbPath)

	txQueryApi = api_gateway.NewTransactionQueryApi(repo)
	txEventListener := api_gateway.NewTransactionEventListener(repo)

	//set mux
	mux := http.NewServeMux()
	httpLogger := kitlog.With(logger, "component", "http")

	err := subscriber.SubscribeTopic("transaction.*", &txEventListener)

	if err != nil {
		panic(err)
	}

	mux.Handle("/", api_gateway.MakeHandler(txQueryApi, httpLogger))
	http.Handle("/", mux)

	go func() {
		log.Println("transport", "http", "address", ipAddress, "msg", "listening")
		errs <- http.ListenAndServe(ipAddress, nil)
	}()

	return nil
}

func initIcode() error {

	log.Println("icode is running...")

	config := conf.GetConfiguration()
	server := rpc.NewServer(config.Engine.Amqp)
	//publisher := pubsub.NewTopicPublisher(config.Engine.Amqp, "Command")

	// service generate
	//commandService := icodeAdapter.NewCommandService(publisher.Publish)

	// git generate
	storeApi := icodeInfra.NewRepositoryService()
	defaultScriptPath := os.Getenv("GOPATH") + "/src/github.com/it-chain/engine/icode/default_setup.sh"
	containerService := icodeService.NewTesseractContainerService(tesseract.Config{
		ShPath: defaultScriptPath,
	})
	api := icodeApi.NewIcodeApi(containerService, storeApi)

	// handler generate
	deployHandler := icodeAdapter.NewDeployCommandHandler(*api)
	unDeployHandler := icodeAdapter.NewUnDeployCommandHandler(*api)
	blockCommandHandler := icodeAdapter.NewBlockCommandHandler(*api)

	server.Register("icode.deploy", deployHandler.HandleDeployCommand)
	server.Register("icode.undeploy", unDeployHandler.HandleUnDeployCommand)
	server.Register("block.execute", blockCommandHandler.HandleBlockExecuteCommand)

	return nil

}

func initPeer() error {
	return nil
}

func initTxPool() error {

	log.Println("txpool is running...")

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)
	server := rpc.NewServer(config.Engine.Amqp)

	//todo get id from pubkey
	tmpPeerID := "tmp peer 1"

	//service
	blockService := txpoolAdapter.NewBlockService(client)
	blockProposalService := txpool.NewBlockProposalService(txQueryApi, blockService, config.Engine.Mode)

	//infra
	txApi := txpoolApi.NewTransactionApi(tmpPeerID)
	txCommandHandler := txpoolAdapter.NewTxCommandHandler(txApi)

	//10초마다 block propose
	txpoolBatch.GetTimeOutBatcherInstance().Run(blockProposalService.ProposeBlock, (time.Duration(config.Txpool.TimeoutMs) * time.Millisecond))

	err := server.Register("transaction.create", txCommandHandler.HandleTxCreateCommand)

	if err != nil {
		panic(err)
	}

	return nil
}

func initConsensus() error {
	return nil
}

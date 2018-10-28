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

package on

import (
	"fmt"
	"os"

	"github.com/it-chain/engine/cmd/on/api-gatewayfx"
	"github.com/it-chain/engine/cmd/on/blockchainfx"
	"github.com/it-chain/engine/cmd/on/commonfx"
	"github.com/it-chain/engine/cmd/on/grpc_gatewayfx"
	"github.com/it-chain/engine/cmd/on/ivmfx"
	"github.com/it-chain/engine/cmd/on/pbftfx"
	"github.com/it-chain/engine/cmd/on/txpoolfx"
	"github.com/it-chain/iLogger"
	"github.com/urfave/cli"
	"go.uber.org/fx"
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

func Action(c *cli.Context) error {
	PrintLogo()
	clearFolder()
	return start()
}

func clearFolder() {

	if err := os.RemoveAll(api_gatewayfx.ApidbPath); err != nil {
		iLogger.Panic(&iLogger.Fields{"err_msg": err.Error()}, "error while clear folder")
		panic(err)
	}

	if err := os.RemoveAll(blockchainfx.BbPath); err != nil {
		iLogger.Panic(&iLogger.Fields{"err_msg": err.Error()}, "error while clear folder")
		panic(err)
	}
	//if err := os.RemoveAll(conf.GetConfiguration().Engine.TmpPath); err != nil {
	//	iLogger.Panic(&iLogger.Fields{"err_msg": err.Error()}, "error while clear folder")
	//	panic(err)
	//}
	//RemoveDirectoryFiles()
}

func start() error {

	app := fx.New(
		commonfx.Module,
		api_gatewayfx.Module,
		grpc_gatewayfx.Module,
		ivmfx.Module,
		txpoolfx.Module,
		blockchainfx.Module,
		pbftfx.Module,
		fx.NopLogger,
	)
	app.Run()
	return nil
}

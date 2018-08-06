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
 *
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/it-chain/engine/icode/mock/handler"
	"github.com/it-chain/sdk"
	"github.com/it-chain/sdk/logger"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Port int `short:"p" long:"port" description:"set port"`
}

func main() {

	logger.EnableFileLogger(true, "./icode.log")
	parser := flags.NewParser(&opts, flags.Default)

	_, err := parser.Parse()

	if err != nil {
		logger.Error(nil, "fail parse args: "+err.Error())
		os.Exit(1)
	}

	fmt.Println("port : " + strconv.Itoa(opts.Port))

	exHandler := &handler.HandlerExample{}
	ibox := sdk.NewIBox(opts.Port)
	ibox.SetHandler(exHandler)
	err = ibox.On(30)

	if err != nil {
		panic(err.Error())
	}
}

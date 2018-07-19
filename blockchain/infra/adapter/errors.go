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

package adapter

import "errors"

var ErrBlockInfoDeliver = errors.New("block info deliver failed")
var ErrGetBlock = errors.New("error when Getting block")
var ErrResponseBlock = errors.New("error when response block")
var ErrGetLastBlock = errors.New("error when get last block")
var ErrSyncCheckResponse = errors.New("error when sync check response")
var ErrEmptyNodeId = errors.New("empty nodeid proposed")
var ErrEmptyBlockSeal = errors.New("empty block seal")
var ErrBlockMissingProperties = errors.New("error when block miss some properties")

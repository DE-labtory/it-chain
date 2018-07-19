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

package blockchain

import "errors"

var ErrTransactionType = errors.New("Wrong transaction type")
var ErrTxListMarshal = errors.New("tx list marshal failed")
var ErrTxListUnmarshal = errors.New("tx list unmarshal failed")
var ErrSetConfig = errors.New("error when get Config")
var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")
var ErrBuildingTxSeal = errors.New("Error in building tx seal")
var ErrBuildingSeal = errors.New("Error in building seal")
var ErrCreatingEvent = errors.New("Error in creating event")
var ErrOnEvent = errors.New("Error on event")

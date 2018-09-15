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

var ErrBlockNil = errors.New("Block nil error")
var ErrBlockTypeCasting = errors.New("Error in type casting block")
var ErrCommandTransactions = errors.New("command's transactions nil or have length of zero")
var ErrCommandSeal = errors.New("command's transactions nil")
var ErrTxHasMissingProperties = errors.New("Tx has missing properties")
var ErrBlockIdNil = errors.New("Error command model ID is nil")
var ErrTxResultsLengthOfZero = errors.New("Error length of tx results is zero")
var ErrTxResultsFail = errors.New("Error not all tx results success")
var ErrCreateEvent = errors.New("Error in creating consent event")

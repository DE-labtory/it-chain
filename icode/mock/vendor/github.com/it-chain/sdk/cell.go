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

package sdk

import "github.com/it-chain/leveldb-wrapper"

type Cell struct {
	DBHandler *leveldbwrapper.DBHandle
}

func NewCell(name string) *Cell {
	path := "./wsdb"
	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
	return &Cell{
		DBHandler: dbProvider.GetDBHandle(name),
	}
}

func (c Cell) PutData(key string, value []byte) error {
	return c.DBHandler.Put([]byte(key), value, true)
}

func (c Cell) GetData(key string) ([]byte, error) {
	value, err := c.DBHandler.Get([]byte(key))
	return value, err
}

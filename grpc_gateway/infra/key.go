/*
 * Copyright 2018 DE-labtory
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

package infra

import (
	"log"

	"github.com/DE-labtory/heimdall/key"
)

func LoadKeyPair(keyPath string, keyType string) (key.PriKey, key.PubKey) {

	km, err := key.NewKeyManager(keyPath)

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err := km.GetKey()

	if err == nil {
		return pri, pub
	}

	pri, pub, err = km.GenerateKey(ConvertToKeyGenOpts(keyType))

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err = km.GetKey()

	return pri, pub
}

func ConvertToKeyGenOpts(keyType string) key.KeyGenOpts {

	switch keyType {
	case "RSA1024":
		return key.RSA1024
	case "RSA2048":
		return key.RSA2048
	case "RSA4096":
		return key.RSA4096
	case "ECDSA256":
		return key.ECDSA256
	default:
		return key.RSA1024
	}
}

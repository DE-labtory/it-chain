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

package common

import (
	"log"

	"errors"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/heimdall"
	"github.com/it-chain/heimdall/config"
	"github.com/it-chain/heimdall/hashing"
	"github.com/it-chain/heimdall/hecdsa"
	"github.com/it-chain/heimdall/keystore"
	"github.com/it-chain/iLogger"
)

func GetNodeID(keyPath string, keyType string) string {
	pri, _ := LoadKeyPair(keyPath, keyType)

	return bifrost.FromPriKey(pri)
}

func LoadKeyPair(keyDirPath string, recoverer bifrost.KeyRecoverer) (heimdall.PriKey, heimdall.PubKey) {

	pri, err := keystore.LoadPriKeyWithoutPwd(keyDirPath, recoverer.(heimdall.KeyRecoverer))
	pub := pri.(heimdall.PriKey).PublicKey()

	if err != nil {
		log.Fatal(err.Error())
	}

	return pri.(heimdall.PriKey), pub
}

type ECDSASigner struct {
	keyDirPath string
	hashOpt    *hashing.HashOpt
}

func (signer *ECDSASigner) Sign(message []byte) ([]byte, error) {
	return hecdsa.SignWithKeyInLocal(signer.keyDirPath, message, signer.hashOpt)
}

type ECDSAVerifier struct {
	signerOpt heimdall.SignerOpts
}

func (verifier *ECDSAVerifier) Verify(peerKey bifrost.Key, signature, message []byte) (bool, error) {
	return hecdsa.Verify(peerKey.(heimdall.PubKey), signature, message, verifier.signerOpt)
}

type ECDSAKeyRecoverer struct {
}

func (rec *ECDSAKeyRecoverer) RecoverKeyFromByte(keyBytes []byte, isPrivate bool) (bifrost.Key, error) {
	recoverer := &hecdsa.KeyRecoverer{}
	key, err := recoverer.RecoverKeyFromByte(keyBytes, isPrivate)
	return key.(bifrost.Key), err
}

func MakeCrypto(secConf *config.Config, keyDirPath string) (bifrost.Crypto, error) {
	signer := &ECDSASigner{
		keyDirPath: keyDirPath,
		hashOpt:    secConf.HashOpt,
	}

	var signerOpt heimdall.SignerOpts
	switch secConf.SigAlgo {
	case "ECDSA":
		signerOpt = hecdsa.NewSignerOpts(secConf.HashOpt)
	case "RSA":
		iLogger.Errorf(nil, "signature algorithm [%s] not supported", secConf.SigAlgo)
		return bifrost.Crypto{}, ErrSigAlgoNotSupported
	default:
		iLogger.Errorf(nil, "signature algorithm [%s] not supported", secConf.SigAlgo)
		return bifrost.Crypto{}, ErrSigAlgoNotSupported
	}

	verifier := &ECDSAVerifier{
		signerOpt: signerOpt,
	}
	keyRecoverer := &ECDSAKeyRecoverer{}

	return bifrost.Crypto{
		Signer:       signer,
		Verifier:     verifier,
		KeyRecoverer: keyRecoverer,
	}, nil

}

func GenerateAndStoreKeyPair(secConf *config.Config, keyDirPath string) (pri heimdall.PriKey, pub heimdall.PubKey, err error) {
	switch secConf.SigAlgo {
	case "ECDSA":
		pri, err = hecdsa.GenerateKey(secConf.KeyGenOpt)
		if err != nil {
			iLogger.Errorf(nil, "key generation error: [%s]", err.Error())
			return nil, nil, ErrKeyGen
		}
		pub = pri.PublicKey()
	case "RSA":
		iLogger.Errorf(nil, "signature algorithm [%s] not supported", secConf.SigAlgo)
		return nil, nil, ErrSigAlgoNotSupported
	default:
		iLogger.Errorf(nil, "signature algorithm [%s] not supported", secConf.SigAlgo)
		return nil, nil, ErrSigAlgoNotSupported
	}

	err = keystore.StorePriKeyWithoutPwd(pri, keyDirPath)
	if err != nil {
		return nil, nil, ErrKeyStore
	}

	return pri, pub, nil
}

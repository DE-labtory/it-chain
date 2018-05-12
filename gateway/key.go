package gateway

import (
	"log"

	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/it-chain-Engine/conf"
)

func loadKeyPair(keyPath string) (key.PriKey, key.PubKey) {

	km, err := key.NewKeyManager(keyPath)

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err := km.GetKey()

	if err == nil {
		return pri, pub
	}

	pri, pub, err = km.GenerateKey(convertToKeyGenOpts(conf.GetConfiguration().Authentication.KeyType))

	if err != nil {
		log.Fatal(err.Error())
	}

	return pri, pub
}

func convertToKeyGenOpts(keyType string) key.KeyGenOpts {

	switch keyType {
	case "RSA1024":
		return key.RSA1024
	case "RSA2048":
		return key.RSA2048
	case "RSA4096":
		return key.RSA4096
	default:
		return key.RSA1024
	}
}

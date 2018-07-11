package infra

import (
	"log"

	"github.com/it-chain/heimdall/key"
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

package auth

import (
	"errors"
)

type keyStore struct {}

func (ks *keyStore) StoreKey(key Key) (err error) {

	if key == nil {
		return errors.New("Failed to get Key Data")
	}

	switch key.(type) {
	case *rsaPrivateKey:
	case *rsaPublicKey:
	case *ecdsaPrivateKey:
	case *ecdsaPublicKey:
	default:
		return errors.New("Unspported Key Type.")
	}
	
	return nil
}

func storePublicKey(path string) (err error) {

}

func storePrivateKey(path string) (err error) {

}



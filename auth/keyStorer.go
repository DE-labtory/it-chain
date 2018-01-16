package auth

import (
	"errors"
)

type keyStorer struct {
	path string
}

func (ks *keyStorer) Store(key Key) (err error) {

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

func storePublicKey(key Key) (err error) {
	return nil
}

func storePrivateKey(key Key) (err error) {
	return nil
}



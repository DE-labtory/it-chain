package auth

import "errors"

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

func publicKeyToPem(pub Key) (pem []byte, err error) {
	return nil, nil
}

func privateKeyToPem(pri Key) (pem []byte, err error) {
	return nil, nil
}

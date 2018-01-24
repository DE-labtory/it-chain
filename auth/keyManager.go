package auth

import (
	"errors"
	"path/filepath"
	"io/ioutil"
	"encoding/hex"
	"os"
)

type keyManager struct {
	path string
}

func (km *keyManager) Store(key Key) (err error) {

	if key == nil {
		return errors.New("Failed to get Key Data")
	}

	switch k := key.(type) {
	case *rsaPrivateKey:
		err = km.storePrivateKey(k)
	case *rsaPublicKey:
		err = km.storePublicKey(k)
	case *ecdsaPrivateKey:
		err = km.storePrivateKey(k)
	case *ecdsaPublicKey:
		err = km.storePublicKey(k)
	default:
		return errors.New("Unspported Key Type.")
	}
	
	return
}

func (km *keyManager) storePublicKey(key Key) (err error) {

	data, err := PublicKeyToPEM(key)
	if err != nil {
		return
	}

	path, err := km.getFullPath(hex.EncodeToString(key.SKI()), "pub")
	if err != nil {
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = ioutil.WriteFile(path, data, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}

func (km *keyManager) storePrivateKey(key Key) (err error) {

	data, err := PrivateKeyToPEM(key)
	if err != nil {
		return
	}

	path, err := km.getFullPath(hex.EncodeToString(key.SKI()), "pri")
	if err != nil {
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = ioutil.WriteFile(path, data, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}

func (km *keyManager) Load(alias string) (key interface{}, err error) {
	return nil, nil
}

func (km *keyManager) loadPublicKey(alias string) (key interface{}, err error) {

}

func (km *keyManager) loadPrivateKey(alias string) (key interface{}, err error) {

}

func (km *keyManager) getSuffix()

func (km *keyManager) getFullPath(alias, suffix string) (path string, err error) {
	if _, err := os.Stat(km.path); os.IsNotExist(err) {
		err = os.MkdirAll(km.path, 0755)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(km.path, alias + "_" + suffix), nil
}



package auth

import (
	"errors"
	"path/filepath"
	"io/ioutil"
	"encoding/hex"
	"os"
	"reflect"
	"strings"
)

type keyManager struct {
	path string

	keyImporters map[reflect.Type]keyImporter
}

func (km *keyManager) Init(path string) {

	if len(path) == 0 {
		km.path = "/KeyRepository"
	} else {
		km.path = path
	}

	keyImporters := make(map[reflect.Type]keyImporter)
	keyImporters[reflect.TypeOf(&RSAPrivateKeyImporterOpts{})] = &rsaPrivateKeyImporter{}
	keyImporters[reflect.TypeOf(&RSAPublicKeyImporterOpts{})] = &rsaPublicKeyImporter{}
	keyImporters[reflect.TypeOf(&ECDSAPrivateKeyImporterOpts{})] = &ecdsaPrivateKeyImporter{}
	keyImporters[reflect.TypeOf(&ECDSAPublicKeyImporterOpts{})] = &ecdsaPublicKeyImporter{}

	km.keyImporters = keyImporters

}

func (km *keyManager) Store(keys... Key) (err error) {

	if len(keys) == 0 {
		return errors.New("Input values should not be NIL")
	}

	for _, key := range keys {
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
	}

	return nil
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

func (km *keyManager) LoadKey() (pub, pri Key, err error) {

	if _, err := os.Stat(km.path); os.IsNotExist(err) {
		return nil, nil, errors.New("Keys are not exist")
	}

	files, err := ioutil.ReadDir(km.path)
	if err != nil {
		return nil, nil, errors.New("Failed to read key repository directory")
	}

	for _, file := range files {

		suffix, valid := km.getSuffix(file.Name())
		if valid == true {
			alias := strings.Split(file.Name(), "_")[0]
			switch suffix {
			case "pri":
				data, err := km.loadPrivateKey(alias)
				if err != nil {
					return nil, nil, err
				}

				pri, err = km.importKey(data)
				if err != nil {
					return nil, nil, err
				}
			case "pub":
				data, err := km.loadPublicKey(alias)
				if err != nil {
					return nil, nil, err
				}

				pri, err = km.importKey(data)
				if err != nil {
					return nil, nil, err
				}
			}
		}
	}

	return

}

func (km *keyManager) loadPrivateKey(alias string) (data interface{}, err error) {

}

func (km *keyManager) loadPublicKey(alias string) (data interface{}, err error) {

}

func (km *keyManager) importKey(data interface{}) (key Key, err error) {

	if data == nil {
		return nil, errors.New("Data should not be NIL")
	}

}

func (km *keyManager) getSuffix(name string) (suffix string, valid bool) {

	if strings.HasSuffix(name, "pri") {
		return "pri", true
	} else if strings.HasSuffix(name, "pub") {
		return "pub", true
	}

	return "", false

}

func (km *keyManager) getFullPath(alias, suffix string) (path string, err error) {
	if _, err := os.Stat(km.path); os.IsNotExist(err) {
		err = os.MkdirAll(km.path, 0755)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(km.path, alias + "_" + suffix), nil
}



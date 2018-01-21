package auth

type rsaPrivateKeyImporter struct {}

func (imp *rsaPrivateKeyImporter) Import(data interface{}) (key Key, err error) {

}

type rsaPublicKeyImporter struct {}

func (imp *rsaPublicKeyImporter) Import(data interface{}) (key Key, err error) {

}

type ecdsaPrivateKeyImporter struct {}

func (imp *ecdsaPrivateKeyImporter) Import(data interface{}) (key Key, err error) {

}

type ecdsaPublicKeyImporter struct {}

func (imp *ecdsaPublicKeyImporter) Import(data interface{}) (key Key, err error) {

}
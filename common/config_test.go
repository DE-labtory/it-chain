package common

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_init_from_file(t *testing.T) {
	assert.Equal(t, "./.leveldb", viper.GetString("database.leveldb.defaultPath"))
}

func TestConfigInit(t *testing.T) {

	var yamlExample = []byte(`
import:
  package: github.com/urfave/cli
  version: ^1.20.0
`)

	err := viper.ReadConfig(bytes.NewBuffer(yamlExample))
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v := viper.GetViper()
	assert.Equal(t, "github.com/urfave/cli", v.GetString("import.package"))
}

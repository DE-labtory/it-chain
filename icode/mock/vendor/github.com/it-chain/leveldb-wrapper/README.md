# leveldb-wrapper

Leveldb-wrapper is a simple library that supporing leveldb and provides an interface for key-value based database.

## Getting Started with leveldb-wrapper

### Installation

```
go get -u github.com/it-chain/leveldb-wrapper
```

### Usage

```Go
package main

import (
	"github.com/it-chain/leveldb-wrapper"
	"fmt"
	"os"
)

func main(){

	path := "./leveldb"
	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
	defer os.RemoveAll(path)

	studentDB := dbProvider.GetDBHandle("Student")
	studentDB.Put([]byte("20164403"),[]byte("JUN"),true)

	name, _ := studentDB.Get([]byte("20164403"))

	fmt.Printf("%s",name)
}
```



## Lincese

Heimdall source code files are made available under the Apache License, Version 2.0 (Apache-2.0), located in the [LICENSE](LICENSE) file.
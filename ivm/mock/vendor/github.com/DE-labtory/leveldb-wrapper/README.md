# leveldb-wrapper

Leveldb-wrapper is a simple library that supporing leveldb and provides an interface for key-value based database.

## KeyValueDB Interface
```
type KeyValueDB interface{
	//OpenDB
	Open()
	//CloseDB
	Close()
	//Get Value by Key
 	Get(key []byte) ([]byte, error)
 	Put(key []byte, value []byte, sync bool) error
 	Delete(key []byte, sync bool) error
 	WriteBatch(KVs map[string][]byte, sync bool) error
	GetIteratorWithPrefix(prefix []byte) KeyValueDBIterator
	GetIterator(startKey []byte, endKey []byte) KeyValueDBIterator
	Snapshot() (map[string][]byte, error)
}

type KeyValueDBIterator interface {
	// It returns whether such pair exist.
	First() bool

	// Last moves the iterator to the last key/value pair. If the iterator
	// only contains one key/value pair then First and Last would moves
	// to the same key/value pair.
	// It returns whether such pair exist.
	Last() bool

	// Seek moves the iterator to the first key/value pair whose key is greater
	// than or equal to the given key.
	// It returns whether such pair exist.
	//
	// It is safe to modify the contents of the argument after Seek returns.
	Seek(key []byte) bool

	// Next moves the iterator to the next key/value pair.
	// It returns whether the iterator is exhausted.
	Next() bool

	// Prev moves the iterator to the previous key/value pair.
	// It returns whether the iterator is exhausted.
	Prev() bool

	// util.Releaser is the interface that wraps basic Release method.
	// When called Release will releases any resources associated with the
	// iterator.
	Release()

	//todo invaild when SetReleaser()
	// util.ReleaseSetter is the interface that wraps the basic SetReleaser
	// method.
	//SetReleaser(asd MyReleaser)

	// TODO: Remove this when ready.
	Valid() bool

	// Error returns any accumulated error. Exhausting all the key/value pairs
	// is not considered to be an error.
	Error() error

	// Key returns the key of the current key/value pair, or nil if done.
	// The caller should not modify the contents of the returned slice, and
	// its contents may change on the next call to any 'seeks method'.
	Key() []byte

	// Value returns the value of the current key/value pair, or nil if done.
	// The caller should not modify the contents of the returned slice, and
	// its contents may change on the next call to any 'seeks method'.
	Value() []byte
}
```

## Getting Started with leveldb-wrapper

### Installation

```
go get -u github.com/DE-labtory/leveldb-wrapper
```

### Usage

```Go
package main

import (
	"github.com/DE-labtory/leveldb-wrapper"
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

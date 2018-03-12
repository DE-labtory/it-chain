package filehelper

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewFileDB(t *testing.T) {
	path := "./test"
	f, err := CreateNewFileDB(path)
	defer f.Remove()
	assert.NoError(t, err)
}

func TestFileDB_Read(t *testing.T) {
	path := "./test"
	f, err := CreateNewFileDB(path)
	defer f.Remove()
	assert.NoError(t, err)

	f.file.Write([]byte("test"))
	r, err := f.Read(0, len([]byte("test")))
	assert.NoError(t, err)
	assert.Equal(t, string(r), "test")
}

func TestFileDB_Write(t *testing.T) {
	path := "./test"
	f, err := CreateNewFileDB(path)
	defer f.Remove()
	assert.NoError(t, err)

	err = f.Write([]byte("test"), true)
	assert.NoError(t, err)
	b := make([]byte, len([]byte("test")))
	_, err = f.file.ReadAt(b, 0)
	assert.Equal(t, string(b), "test")
}
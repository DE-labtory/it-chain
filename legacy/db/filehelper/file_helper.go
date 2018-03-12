package filehelper

import (
	"os"
)

type fileDB struct {
	filePath string
	file     *os.File
}

func CreateNewFileDB(filePath string) (*fileDB, error) {
	f := &fileDB{filePath: filePath}
	return f, f.Open()
}

func (f *fileDB) Open() error {
	file, err := os.OpenFile(f.filePath, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	f.file = file
	return nil
}

func (f *fileDB) Close() error {
	return f.file.Close()
}

func (f *fileDB) Remove() error {
	err := f.Close()
	if err != nil {
		return err
	}
	return os.Remove(f.filePath)
}

func (f *fileDB) Read(offset int, length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := f.file.ReadAt(b, int64(offset))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (f *fileDB) Write(b []byte, sync bool) error {
	_, err := f.file.Write(b)
	if err != nil {
		return err
	}
	if sync {
		return f.file.Sync()
	}
	return nil
}
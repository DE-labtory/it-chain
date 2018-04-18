package leveldb

import (
	"os"
	"testing"

	"github.com/it-chain/it-chain-Engine/icode/domain/model"
	"github.com/stretchr/testify/assert"
)

func newRepo() *ICodeMetaRepository {
	tmpPath := "./tmp"
	defer os.RemoveAll(tmpPath)

	iCodeRepo := NewICodeMetaRepository(tmpPath)

	return iCodeRepo
}

func TestNewICodeMetaRepositoryst(t *testing.T) {

	//given
	tmpPath := "./tmp"
	defer os.RemoveAll(tmpPath)

	//when
	iCodeRepo := NewICodeMetaRepository(tmpPath)

	defer iCodeRepo.Close()

	//then
	assert.DirExists(t, tmpPath)
}

func TestICodeMetaRepository_Save(t *testing.T) {

	//given
	iCodeRepo := newRepo()
	defer iCodeRepo.Close()

	icodeMeta := model.NewICodeMeta("test", "est", "Test", "commit")

	//when
	err := iCodeRepo.Save(*icodeMeta)
	assert.NoError(t, err)

	//then
	b, err := iCodeRepo.leveldb.Get([]byte(icodeMeta.ID))
	ficodeMeta := &model.ICodeMeta{}
	err = Deserialize(b, ficodeMeta)
	assert.NoError(t, err)

	assert.Equal(t, icodeMeta, ficodeMeta)
}

func TestICodeMetaRepository_FindByID(t *testing.T) {

	//given
	iCodeRepo := newRepo()
	defer iCodeRepo.Close()

	icodeMeta := model.NewICodeMeta("test", "est", "Test", "commit")

	err := iCodeRepo.Save(*icodeMeta)
	assert.NoError(t, err)

	//when
	ficodeMeta, err := iCodeRepo.FindByID(icodeMeta.ID)
	assert.NoError(t, err)

	//then
	assert.Equal(t, icodeMeta, ficodeMeta)
}

func TestICodeMetaRepository_FindAll(t *testing.T) {

	//given
	iCodeRepo := newRepo()
	defer iCodeRepo.Close()
	icodeMeta1 := model.NewICodeMeta("test1", "est", "Test", "commit")
	icodeMeta2 := model.NewICodeMeta("test2", "est", "Test", "commit")
	icodeMeta3 := model.NewICodeMeta("test3", "est", "Test", "commit")

	err := iCodeRepo.Save(*icodeMeta1)
	assert.NoError(t, err)

	err = iCodeRepo.Save(*icodeMeta2)
	assert.NoError(t, err)

	err = iCodeRepo.Save(*icodeMeta3)
	assert.NoError(t, err)

	//when
	ficodeMetas, err := iCodeRepo.FindAll()
	assert.NoError(t, err)

	//then
	assert.Contains(t, ficodeMetas, icodeMeta1)
	assert.Contains(t, ficodeMetas, icodeMeta2)
	assert.Contains(t, ficodeMetas, icodeMeta3)
}

func TestICodeMetaRepository_Remove(t *testing.T) {

	//given
	iCodeRepo := newRepo()
	defer iCodeRepo.Close()

	icodeMeta1 := model.NewICodeMeta("test1", "est", "Test", "commit")

	err := iCodeRepo.Save(*icodeMeta1)
	assert.NoError(t, err)

	//when
	err = iCodeRepo.Remove(icodeMeta1.ID)
	assert.NoError(t, err)

	//then
	fiCodeMeta, err := iCodeRepo.FindByID(icodeMeta1.ID)
	assert.NoError(t, err)
	assert.Nil(t, fiCodeMeta)
}

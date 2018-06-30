package blockchain

//func TestCreateGenesisBlock(t *testing.T) {
//
//	GenesisFilePath := "./GenesisBlockConfig.json"
//	defer os.Remove(GenesisFilePath)
//	wrongFilePath := "./WrongFileName.json"
//	GenesisBlockConfigJson := []byte(`{
//								  "Seal":[],
//								  "PrevSeal":[],
//								  "Height":0,
//								  "TxList":[],
//								  "TxSeal":[],
//								  "TimeStamp":"0001-01-01T00:00:00-00:00",
//								  "Creator":[]
//								}`)
//	validator := new(impl.DefaultValidator)
//	var tempBlock impl.DefaultBlock
//	_ = json.Unmarshal(GenesisBlockConfigJson, &tempBlock)
//	GenesisBlockConfigByte, _ := json.Marshal(tempBlock)
//	_ = ioutil.WriteFile(GenesisFilePath, GenesisBlockConfigByte, 0644)
//
//	GenesisBlock, err1 := CreateGenesisBlock(GenesisFilePath)
//	expectedSeal, _ := validator.BuildSeal(GenesisBlock)
//	assert.NoError(t, err1)
//	assert.Equal(t, expectedSeal, GenesisBlock.Seal)
//	assert.Equal(t, make([]byte, 0), GenesisBlock.PrevSeal)
//	assert.Equal(t, uint64(0), GenesisBlock.Height)
//	assert.Equal(t, make([]*impl.DefaultTransaction, 0), GenesisBlock.TxList)
//	assert.Equal(t, make([][]byte, 0), GenesisBlock.TxSeal)
//	assert.Equal(t, time.Now().String()[:19], GenesisBlock.Timestamp.String()[:19])
//	assert.Equal(t, make([]byte, 0), GenesisBlock.Creator)
//
//	_, err2 := CreateGenesisBlock(wrongFilePath)
//	assert.Error(t, err2)
//}

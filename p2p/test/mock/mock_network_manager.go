package mock

type MockNetworkManager struct {
	mockProcessTable map[string]MockProcess
}

func (mnm *MockNetworkManager) AddMockProcess(mockProcess MockProcess) {

	mnm.mockProcessTable[mockProcess.GetId()] = mockProcess
}

func (mnm *MockNetworkManager) FindMockProcess(id string) MockProcess {

	return mnm.mockProcessTable[id]
}

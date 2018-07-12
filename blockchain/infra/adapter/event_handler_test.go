package adapter_test

//todo eventstore를 활용한 testcase재 작성필요
//func TestEventHandler_HandleBlockAddToPoolEvent(t *testing.T) {
//	tests := map[string]struct {
//		input struct {
//			blockchain.BlockAddToPoolEvent
//		}
//		err error
//	}{
//		"success": {
//			input: struct {
//				blockchain.BlockAddToPoolEvent
//			}{BlockAddToPoolEvent: blockchain.BlockAddToPoolEvent{
//				Block: &blockchain.DefaultBlock{
//					Height: uint64(12),
//				},
//			}},
//			err: nil,
//		},
//		"block nil test": {
//			input: struct {
//				blockchain.BlockAddToPoolEvent
//			}{BlockAddToPoolEvent: blockchain.BlockAddToPoolEvent{
//				Block: nil,
//			}},
//			err: adapter.ErrBlockNil,
//		},
//	}
//
//	// When
//	nr := MockPeerRepository{}
//	er := MockEventRepository{}
//	rp := adapter.RepositoryProjector{
//		PeerRepository:  nr,
//		EventRepository: er,
//	}
//
//	blockApi := MockEventListenerBlockApi{}
//
//	eventHandler := adapter.NewEventHandler(rp, blockApi)
//
//	for testName, test := range tests {
//		t.Logf("running test case %s", testName)
//
//		// When
//		err := eventHandler.HandleBlockAddToPoolEvent(test.input.BlockAddToPoolEvent)
//
//		// Then
//		assert.Equal(t, err, test.err)
//	}
//}

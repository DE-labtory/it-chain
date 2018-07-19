package adapter

type EventHandler struct {
	blockApi BlockApi
}

//ToDo: 미완성 부분에 대한 주석처리
//
//func NewEventHandler(api BlockApi) *EventHandler {
//	return &EventHandler{
//		blockApi: api,
//	}
//}
//
//// TODO: write test case
//func (eh *EventHandler) HandleBlockAddToPoolEvent(event blockchain.BlockAddToPoolEvent) error {
//	if err := isBlockHasMissingProperty(event); err != nil {
//		return err
//	}
//	height := event.Height
//	err := eh.blockApi.CheckAndSaveBlockFromPool(height)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func isBlockHasMissingProperty(event blockchain.BlockAddToPoolEvent) error {
//	if event.Seal == nil || event.PrevSeal == nil || event.Height == 0 ||
//		event.TxList == nil || event.TxSeal == nil || event.Timestamp.IsZero() || event.Creator == nil {
//		return ErrBlockMissingProperties
//	}
//	return nil
//}

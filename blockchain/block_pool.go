package blockchain

// ToDo: 삭제 제안 - junk_sound
//type BlockPool interface {
//	Add(block Block) error
//	Get(height BlockHeight) Block
//	Delete(height Block)
//}
//
//var BLOCK_POOL_AID = "BLOCK_POOL_AID"
//
//type BlockPoolModel struct {
//	midgard.AggregateModel
//	Pool map[BlockHeight]Block
//}
//
//func NewBlockPool() *BlockPoolModel {
//	return &BlockPoolModel{
//		AggregateModel: midgard.AggregateModel{
//			ID: BLOCK_POOL_AID,
//		},
//		Pool: make(map[BlockHeight]Block),
//	}
//}
//
//func (p *BlockPoolModel) Add(block Block) error {
//	event, err := createBlockAddToPoolEvent(block)
//	if err != nil {
//		return err
//	}
//
//	eventstore.Save(BLOCK_POOL_AID, event)
//
//	p.On(&event)
//
//	return nil
//}
//
//func (p *BlockPoolModel) Get(height BlockHeight) Block {
//	return p.Pool[height]
//}
//
//func (p *BlockPoolModel) Delete(block Block) {
//	event := createBlockRemoveFromPoolEvent(block)
//	eventstore.Save(BLOCK_POOL_AID, event)
//
//	p.On(&event)
//}
//
//func createBlockAddToPoolEvent(block Block) (BlockAddToPoolEvent, error) {
//	txListBytes, err := common.Serialize(block.GetTxList())
//	if err != nil {
//		return BlockAddToPoolEvent{}, ErrTxListMarshal
//	}
//
//	return BlockAddToPoolEvent{
//		EventModel: midgard.EventModel{
//			ID: BLOCK_POOL_AID,
//		},
//		Seal:      block.GetSeal(),
//		PrevSeal:  block.GetPrevSeal(),
//		Height:    block.GetHeight(),
//		TxList:    txListBytes,
//		TxSeal:    block.GetTxSeal(),
//		Timestamp: block.GetTimestamp(),
//		Creator:   block.GetCreator(),
//	}, nil
//}
//
//func createBlockRemoveFromPoolEvent(block Block) BlockRemoveFromPoolEvent {
//	return BlockRemoveFromPoolEvent{
//		EventModel: midgard.EventModel{
//			ID: BLOCK_POOL_AID,
//		},
//		Height: block.GetHeight(),
//	}
//}
//
//func (p *BlockPoolModel) GetID() string {
//	return BLOCK_POOL_AID
//}
//
//func (p *BlockPoolModel) On(event midgard.Event) error {
//	switch v := event.(type) {
//
//	case *BlockAddToPoolEvent:
//		block, err := createBlockFromAddToPoolEvent(v)
//		if err != nil {
//			return err
//		}
//		(p.Pool)[v.Height] = block
//
//	case *BlockRemoveFromPoolEvent:
//		delete(p.Pool, v.Height)
//
//	default:
//		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
//	}
//	return nil
//}
//
//func createBlockFromAddToPoolEvent(event *BlockAddToPoolEvent) (Block, error) {
//	txList, err := deserializeTxList(event.TxList)
//	if err != nil {
//		return &DefaultBlock{}, ErrTxListUnmarshal
//	}
//
//	return &DefaultBlock{
//		Seal:      event.Seal,
//		PrevSeal:  event.PrevSeal,
//		Height:    event.Height,
//		TxList:    txList,
//		TxSeal:    event.TxSeal,
//		Timestamp: event.Timestamp,
//		Creator:   event.Creator,
//	}, nil
//}

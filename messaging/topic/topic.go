package topic

type Event int

const (
	MessageCreated Event = iota
	ConsensusMessagePublishEvent
	BlockConfirmEvent
	//txpool Event
	TransactionReceiveEvent
	TransactionSendEvent
	BlockProposeEvent
)

func (e Event) String() string {
	switch e {
	case MessageCreated:
		return "MessageCreated"
	case ConsensusMessagePublishEvent:
		return "ConsensusMessagePublishEvent"
	case BlockConfirmEvent:
		return "BlockConfirmEvent"
	case TransactionReceiveEvent:
		return "TransactionReceiveEvent"
	case TransactionSendEvent:
		return "TransactionSendEvent"
	case BlockProposeEvent:
		return "BlockProposeEvent"
	}

	return "error"
}

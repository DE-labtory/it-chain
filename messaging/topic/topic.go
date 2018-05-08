package topic

type Event int

const (
	MessageCreated Event = iota
	ConsensusMessagePublishEvent
	BlockConfirmEvent
<<<<<<< HEAD
	MessageDeliverEvent
	NewConnEvent
=======
  ConnectionCreated
	//txpool Event
	TransactionReceiveEvent
	TransactionSendEvent
	BlockProposeEvent
>>>>>>> 251a47ac18929415ccdda952ddf228bf5ad7077c
)

func (e Event) String() string {
	switch e {
	case MessageCreated:
		return "MessageCreated"
	case ConsensusMessagePublishEvent:
		return "ConsensusMessagePublishEvent"
	case BlockConfirmEvent:
		return "BlockConfirmEvent"
<<<<<<< HEAD
	case MessageDeliverEvent:
		return "MessageDeliverEvent"
	case NewConnEvent:
		return "NewConnEvent"
=======
	case TransactionReceiveEvent:
		return "TransactionReceiveEvent"
	case TransactionSendEvent:
		return "TransactionSendEvent"
	case BlockProposeEvent:
		return "BlockProposeEvent"
	case ConnectionCreated:
		return "ConnectionCreated"
>>>>>>> 251a47ac18929415ccdda952ddf228bf5ad7077c
	}

	return "error"
}

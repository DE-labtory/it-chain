package topic

type Event int

const (
	MessageCreated Event = iota
	ConsensusMessagePublishEvent
	BlockConfirmEvent
)

func (e Event) String() string {
	switch e {
	case MessageCreated:
		return "MessageCreated"
	case ConsensusMessagePublishEvent:
		return "ConsensusMessagePublishEvent"
	case BlockConfirmEvent:
		return "BlockConfirmEvent"
	}

	return "error"
}

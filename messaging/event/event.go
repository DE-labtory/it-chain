package event

type Event int

const (
	MessageCreated Event = iota
)

func (e Event) String() string {
	switch e {
	case MessageCreated:
		return "MessageCreated"
	}

	return "error"
}
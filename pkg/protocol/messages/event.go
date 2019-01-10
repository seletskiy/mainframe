package messages

type Event struct {
	Tick int64
	Kind string
}

func (*Event) Tag() string {
	return "event"
}

func (message *Event) Serialize() []Arg {
	return []Arg{
		{"tick", message.Tick},
		{"kind", message.Kind},
	}
}

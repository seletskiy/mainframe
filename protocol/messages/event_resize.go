package messages

type EventResize struct {
	Event

	Width  int
	Height int

	Columns int
	Rows    int
}

func (message *EventResize) Serialize() []Arg {
	return append(
		message.Event.Serialize(),
		Arg{"columns", message.Columns},
		Arg{"rows", message.Rows},
		Arg{"width", message.Width},
		Arg{"height", message.Height},
	)
}

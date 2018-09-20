package messages

type EventInput struct {
	Event

	Char rune

	Ctrl  bool
	Shift bool
	Alt   bool
	Super bool
}

func (message *EventInput) Serialize() []Arg {
	return append(
		message.Event.Serialize(),
		Arg{"char", string(message.Char)},
		Arg{"shift", message.Shift},
		Arg{"ctrl", message.Ctrl},
		Arg{"alt", message.Alt},
		Arg{"super", message.Super},
	)
}

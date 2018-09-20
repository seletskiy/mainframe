package messages

type EventKeyboard struct {
	Event

	Symbol string
	Code   int

	Press   bool
	Release bool
	Repeat  bool

	Ctrl  bool
	Shift bool
	Alt   bool
	Super bool
}

func (message *EventKeyboard) Serialize() []Arg {
	return append(
		message.Event.Serialize(),
		Arg{"symbol", message.Symbol},
		Arg{"press", message.Press},
		Arg{"release", message.Release},
		Arg{"repeat", message.Repeat},
		Arg{"code", message.Code},
		Arg{"shift", message.Shift},
		Arg{"ctrl", message.Ctrl},
		Arg{"alt", message.Alt},
		Arg{"super", message.Super},
	)
}

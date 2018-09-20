package messages

type Error struct {
	Message string
}

func (*Error) Tag() string {
	return "error"
}

func (message *Error) Serialize() []Arg {
	return []Arg{
		{"message", message.Message},
	}
}

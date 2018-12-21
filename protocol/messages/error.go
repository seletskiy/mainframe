package messages

type Error struct {
	Message string
}

func (error *Error) Error() string {
	return error.Message
}

func (*Error) Tag() string {
	return "error"
}

func (message *Error) Serialize() []Arg {
	return []Arg{
		{"message", message.Message},
	}
}

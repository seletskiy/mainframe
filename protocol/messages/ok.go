package messages

type OK struct {
	Args []Arg
}

func (*OK) Tag() string {
	return "ok"
}

func (message *OK) Serialize() []Arg {
	return message.Args
}

func (message *OK) Set(name string, value interface{}) *OK {
	message.Args = append(message.Args, Arg{name, value})

	return message
}

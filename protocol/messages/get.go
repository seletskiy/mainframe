package messages

type Get struct {
	Font struct {
		Set bool
	}
}

func (message *Get) Tag() string {
	return "get"
}

func (message *Get) Serialize() []Arg {
	args := []Arg{}

	if message.Font.Set {
		args = append(args, Arg{"font", true})
	}

	return args
}

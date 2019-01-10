package messages

type Reshape struct {
	*Size

	X *int
	Y *int
}

func (message *Reshape) Tag() string {
	return "reshape"
}

func (message *Reshape) Serialize() []Arg {
	args := []Arg{
		{"x", message.X},
		{"y", message.Y},
	}

	if message.Size != nil {
		args = append(args, message.Size.Serialize()...)
	}

	return args
}

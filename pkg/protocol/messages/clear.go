package messages

type Clear struct {
	X *int
	Y *int

	Rows    *int
	Columns *int
}

func (*Clear) Tag() string {
	return "clear"
}

func (message *Clear) Serialize() []Arg {
	return []Arg{
		{"x", message.X},
		{"y", message.Y},
		{"rows", message.Rows},
		{"columns", message.Columns},
	}
}

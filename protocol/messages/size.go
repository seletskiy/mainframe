package messages

type Size struct {
	Width  *int
	Height *int

	Rows    *int
	Columns *int
}

func (message *Size) Serialize() []Arg {
	return []Arg{
		{"width", message.Width},
		{"height", message.Height},
		{"rows", message.Rows},
		{"columns", message.Columns},
	}
}

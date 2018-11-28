package messages

type Open struct {
	Width  *int
	Height *int

	Rows    *int
	Columns *int

	X *int
	Y *int

	Title string

	Raw      bool
	Hidden   bool
	Fixed    bool
	Bare     bool
	Floating bool
}

func (Open) Tag() string {
	return "open"
}

func (message *Open) Serialize() []Arg {
	return []Arg{
		{"width", message.Width},
		{"height", message.Height},
		{"rows", message.Rows},
		{"columns", message.Columns},
		{"x", message.X},
		{"y", message.Y},
		{"title", message.Title},
		{"raw", message.Raw},
		{"hidden", message.Hidden},
		{"fixed", message.Fixed},
		{"bare", message.Bare},
		{"floating", message.Floating},
	}
}

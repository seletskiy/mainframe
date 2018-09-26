package messages

type Open struct {
	Width  *int
	Height *int

	X *int
	Y *int

	Title string

	Raw       bool
	Hidden    bool
	Fixed     bool
	Decorated bool
	Floating  bool
}

func (Open) Tag() string {
	return "open"
}

func (message *Open) Serialize() []Arg {
	return []Arg{
		{"width", message.Width},
		{"height", message.Height},
		{"x", message.X},
		{"y", message.Y},
		{"title", message.Title},
		{"raw", message.Raw},
		{"hidden", message.Hidden},
		{"fixed", message.Fixed},
		{"decorated", message.Decorated},
		{"floating", message.Floating},
	}
}

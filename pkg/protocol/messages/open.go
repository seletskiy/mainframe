package messages

type Open struct {
	*Size

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
	args := []Arg{
		{"x", message.X},
		{"y", message.Y},
		{"title", message.Title},
		{"raw", message.Raw},
		{"hidden", message.Hidden},
		{"fixed", message.Fixed},
		{"bare", message.Bare},
		{"floating", message.Floating},
	}

	if message.Size != nil {
		args = append(args, message.Size.Serialize()...)
	}

	return args
}

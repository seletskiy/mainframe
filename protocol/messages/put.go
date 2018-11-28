package messages

import (
	"image/color"
)

type Put struct {
	X int
	Y int

	Columns *int
	Rows    *int

	Foreground *color.RGBA
	Background *color.RGBA
	Text       *string

	Tick      *int
	Exclusive *bool
}

func (message *Put) Tag() string {
	return "put"
}

func (message *Put) Serialize() []Arg {
	return []Arg{
		{"x", message.X},
		{"y", message.Y},
		{"columns", message.Columns},
		{"rows", message.Rows},
		{"fg", message.Foreground},
		{"bg", message.Background},
		{"text", message.Text},
		{"tick", message.Tick},
		{"exclusive", message.Exclusive},
	}
}

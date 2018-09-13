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

func (put *Put) Tag() string {
	return "put"
}

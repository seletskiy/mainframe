package messages

import (
	"image/color"
)

type Put struct {
	X int
	Y int

	Width  *int
	Height *int

	Foreground *color.RGBA
	Background *color.RGBA
	Text       *string

	Tick      *int
	Exclusive *bool
}

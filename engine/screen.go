package engine

import (
	"image/color"
	"sync"

	"github.com/seletskiy/mainframe/fonts"
	"github.com/seletskiy/mainframe/protocol/messages"
)

const (
	AttrEmpty      = 0
	AttrGlyph      = 1
	AttrForeground = 2
	AttrBackground = 4
)

type Screen struct {
	font *fonts.Font

	width  int
	height int
	cells  []int32
	attrs  []int32
	colors []int32

	lock sync.Mutex
}

func NewScreen(width, height int, font *fonts.Font) *Screen {
	return &Screen{
		width:  width / font.Meta.Width,
		height: height / font.Meta.Height,
		font:   font,
		cells:  make([]int32, 2*width*height),
		attrs:  make([]int32, width*height),
		colors: make([]int32, 2*width*height),
	}
}

func (screen *Screen) SetGlyph(x, y int, char string) bool {
	if x > screen.width {
		return false
	}

	if y > screen.height {
		return false
	}

	glyph := screen.font.Glyphs[char]
	if glyph == nil {
		// TODO: draw some missing char
		return true
	}

	pos := x + y*screen.width

	screen.cells[pos*2] = int32(glyph.Column)
	screen.cells[pos*2+1] = int32(glyph.Row)
	screen.attrs[pos] |= AttrGlyph

	return true
}

func (screen *Screen) SetForeground(x, y int, fg *color.RGBA) bool {
	if x > screen.width {
		return false
	}

	if y > screen.height {
		return false
	}

	pos := x + y*screen.width

	screen.colors[pos*2] = int32(fg.R)<<16 + int32(fg.G)<<8 + int32(fg.B)
	screen.attrs[pos] |= AttrForeground

	return true
}

func (screen *Screen) SetBackground(x, y int, bg *color.RGBA) bool {
	if x > screen.width {
		return false
	}

	if y > screen.height {
		return false
	}

	pos := x + y*screen.width

	screen.colors[pos*2+1] = int32(bg.R)<<16 + int32(bg.G)<<8 + int32(bg.B)
	screen.attrs[pos] |= AttrBackground

	return true
}

func (screen *Screen) Put(message *messages.Put) bool {
	var width int
	var height int

	if message.Width != nil {
		width = *message.Width
	}

	if message.Height != nil {
		height = *message.Height
	} else {
		height = 1
	}

	if message.Text != nil {
		text := *message.Text

		if message.Width == nil {
			width = len(text)
		}

		i := 0

	text:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if i >= len(text) {
					break text
				}

				char := string(text[i])
				if !screen.SetGlyph(x+message.X, y+message.Y, char) {
					break text
				}

				i++
			}
		}
	}

	if message.Foreground != nil {
	foreground:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if !screen.SetForeground(
					x+message.X,
					y+message.Y,
					message.Foreground,
				) {
					break foreground
				}
			}
		}
	}

	if message.Background != nil {
	background:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if !screen.SetBackground(
					x+message.X,
					y+message.Y,
					message.Background,
				) {
					break background
				}
			}
		}
	}

	return true
}

func (screen *Screen) GetSize() int {
	return screen.width * screen.height
}

func (screen *Screen) Resize(width, height int) {
	width /= screen.font.Meta.Width
	height /= screen.font.Meta.Height

	cells := make([]int32, 2*width*height)
	attrs := make([]int32, width*height)

	for y := 0; y < screen.height; y++ {
		if y >= height {
			break
		}

		for x := 0; x < screen.width; x++ {
			if x >= width {
				break
			}

			var (
				from = x + y*screen.width
				to   = x + y*width
			)

			cells[to*2] = screen.cells[from*2]
			cells[to*2+1] = screen.cells[from*2+1]
			attrs[to] = screen.attrs[from]
		}
	}

	screen.width = width
	screen.height = height
	screen.cells = cells
	screen.attrs = attrs
}

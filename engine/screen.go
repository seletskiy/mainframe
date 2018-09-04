package engine

import (
	"sync"

	"github.com/seletskiy/mainframe/fonts"
)

type Screen struct {
	font *fonts.Font

	width  int
	height int
	cells  []int32
	attrs  []int32

	lock sync.Mutex
}

func NewScreen(width, height int, font *fonts.Font) *Screen {
	return &Screen{
		width:  width / font.Meta.Width,
		height: height / font.Meta.Height,
		font:   font,
		cells:  make([]int32, 2*width*height),
		attrs:  make([]int32, width*height),
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
	screen.attrs[pos] = 1

	return true
}

func (screen *Screen) SetText(x, y int, text string) bool {
	for _, char := range text {
		if !screen.SetGlyph(x, y, string(char)) {
			return false
		}

		x += 1
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

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

	width   int
	height  int
	rows    int
	columns int

	cells  []int32
	attrs  []int32
	colors []int32

	lock sync.Mutex
}

func NewScreen(width, height int, font *fonts.Font) *Screen {
	var (
		columns = width / font.Meta.Width
		rows    = height / font.Meta.Height
	)

	return &Screen{
		width:   width,
		height:  height,
		columns: columns,
		rows:    rows,
		font:    font,
		cells:   make([]int32, 2*rows*columns),
		attrs:   make([]int32, rows*columns),
		colors:  make([]int32, 2*rows*columns),
	}
}

func (screen *Screen) SetGlyph(x, y int, char string) bool {
	if x > screen.columns {
		return false
	}

	if y > screen.rows {
		return false
	}

	glyph := screen.font.Glyphs[char]
	if glyph == nil {
		// TODO: draw some missing char
		return true
	}

	pos := x + y*screen.columns

	screen.cells[pos*2] = int32(glyph.Column)
	screen.cells[pos*2+1] = int32(glyph.Row)
	screen.attrs[pos] |= AttrGlyph

	return true
}

func (screen *Screen) SetForeground(x, y int, fg *color.RGBA) bool {
	if x >= screen.columns {
		return false
	}

	if y >= screen.rows {
		return false
	}

	pos := x + y*screen.columns

	screen.colors[pos*2] = int32(fg.R)<<16 + int32(fg.G)<<8 + int32(fg.B)
	screen.attrs[pos] |= AttrForeground

	return true
}

func (screen *Screen) SetBackground(x, y int, bg *color.RGBA) bool {
	if x >= screen.columns {
		return false
	}

	if y >= screen.rows {
		return false
	}

	pos := x + y*screen.columns

	screen.colors[pos*2+1] = int32(bg.R)<<16 + int32(bg.G)<<8 + int32(bg.B)
	screen.attrs[pos] |= AttrBackground

	return true
}

func (screen *Screen) Put(message *messages.Put) bool {
	screen.lock.Lock()
	defer screen.lock.Unlock()

	var columns int
	var rows int

	if message.Columns != nil {
		columns = *message.Columns
	}

	if message.Rows != nil {
		rows = *message.Rows
	} else {
		rows = 1
	}

	if message.Text != nil {
		text := *message.Text

		if message.Columns == nil {
			columns = len(text)
		}

		i := 0

	text:
		for y := 0; y < rows; y++ {
			for x := 0; x < columns; x++ {
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
		for y := 0; y < rows; y++ {
			for x := 0; x < columns; x++ {
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
		for y := 0; y < rows; y++ {
			for x := 0; x < columns; x++ {
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
	return screen.columns * screen.rows
}

func (screen *Screen) Resize(width, height int) {
	screen.lock.Lock()
	defer screen.lock.Unlock()

	var (
		columns = width / screen.font.Meta.Width
		rows    = height / screen.font.Meta.Height

		cells  = make([]int32, 2*rows*columns)
		attrs  = make([]int32, rows*columns)
		colors = make([]int32, 2*rows*columns)
	)

	for y := 0; y < screen.rows; y++ {
		if y >= rows {
			break
		}

		for x := 0; x < screen.columns; x++ {
			if x >= columns {
				break
			}

			var (
				from = x + y*screen.columns
				to   = x + y*columns
			)

			cells[to*2] = screen.cells[from*2]
			cells[to*2+1] = screen.cells[from*2+1]
			attrs[to] = screen.attrs[from]
			colors[to*2] = screen.colors[from*2]
			colors[to*2+1] = screen.colors[from*2+1]
		}
	}

	screen.width = width
	screen.height = height
	screen.rows = rows
	screen.columns = columns
	screen.cells = cells
	screen.attrs = attrs
	screen.colors = colors
}

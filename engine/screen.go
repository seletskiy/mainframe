package engine

import (
	"image/color"
	"strconv"
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

	regions map[string]ScreenRegion

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

		regions: make(map[string]ScreenRegion),
	}
}

func (screen *Screen) SetGlyph(x, y int, char string) bool {
	if x >= screen.columns {
		return false
	}

	if y >= screen.rows {
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

	var (
		address = screen.getRegionID(message.X, message.Y)
		region  = screen.regions[address]
	)

	if region.Exclusive {
		screen.Clear(message.X, message.Y, region.Rows, region.Columns)

		region.Exclusive = false
	}

	var (
		offscreen bool
		columns   int
		rows      int
	)

	if message.Columns != nil {
		columns = *message.Columns
	}

	if message.Rows != nil {
		rows = *message.Rows
	} else {
		rows = 1
	}

	switch {
	case message.Exclusive:
		region.Exclusive = true
	default:
		delete(screen.regions, address)
	}

	if message.Text != nil {
		text := *message.Text

		if message.Columns == nil {
			columns = len(text)
		}

		var i int

		for _, char := range text {
			x := i % columns
			y := i / columns
			if y > rows {
				break
			}

			if char == '\n' {
				i += (columns - i%columns)
				continue
			}

			x += message.X
			y += message.Y

			if !screen.SetGlyph(x, y, string(char)) {
				offscreen = true
			} else {
				if y >= region.Rows {
					region.Rows = y + 1
				}

				if x >= region.Columns {
					region.Columns = x + 1
				}
			}

			i++
		}
	}

	screen.regions[address] = region

	if message.Foreground != nil {
		for y := 0; y < rows; y++ {
			for x := 0; x < columns; x++ {
				if !screen.SetForeground(
					x+message.X,
					y+message.Y,
					message.Foreground,
				) {
					offscreen = true
				}
			}
		}
	}

	if message.Background != nil {
		for y := 0; y < rows; y++ {
			for x := 0; x < columns; x++ {
				if !screen.SetBackground(
					x+message.X,
					y+message.Y,
					message.Background,
				) {
					offscreen = true
				}
			}
		}
	}

	return !offscreen
}

func (screen *Screen) GetSize() int {
	return screen.columns * screen.rows
}

func (screen *Screen) Resize(width, height int) (int, int) {
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

	return columns, rows
}

func (screen *Screen) Clear(x int, y int, rows int, columns int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			screen.attrs[(x+j)+(y+i)*screen.columns] = AttrEmpty
		}
	}
}

func (screen *Screen) getRegionID(x int, y int) string {
	return strconv.Itoa(x) + ":" + strconv.Itoa(y)
}

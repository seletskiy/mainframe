package fonts

import (
	"image"
	"image/draw"
	"io/ioutil"
	"math"

	"golang.org/x/image/math/fixed"

	"github.com/reconquest/karma-go"
	"github.com/seletskiy/freetype"
	"github.com/seletskiy/freetype/truetype"

	xfont "golang.org/x/image/font"
)

type Glyph struct {
	Row    int
	Column int
	Char   string
}

type Font struct {
	Image  *image.RGBA
	Glyphs map[string]*Glyph

	handle  *truetype.Font
	metrics struct {
		length int
		width  int
		height int

		descender int
	}
}

type FontDPI float64
type FontSize float64
type FontHinting bool

func Load(name string, opts ...interface{}) (*Font, error) {
	font := &Font{
		Glyphs: make(map[string]*Glyph),
	}

	err := font.load(name)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to load font file",
		)
	}

	var (
		dpi     float64
		size    float64
		hinting xfont.Hinting
	)

	for _, opt := range opts {
		switch opt := opt.(type) {
		case FontDPI:
			dpi = float64(opt)
		case FontSize:
			size = float64(opt)
		case FontHinting:
			if opt {
				hinting = xfont.HintingFull
			}
		}
	}

	err = font.raster(dpi, size, hinting)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to rasterize font",
		)
	}

	return font, nil
}

func (font *Font) GetWidth() int {
	return font.metrics.width
}

func (font *Font) GetHeight() int {
	return font.metrics.height
}

func (font *Font) load(path string) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return karma.Format(
			err,
			"unable to read font file",
		)
	}

	font.handle, err = truetype.Parse(body)
	if err != nil {
		return karma.Format(
			err,
			"unable to parse font",
		)
	}

	return nil
}

// raster converts given vector font into rasterized image that contains all
// characters defined in font.
//
// Rasterizer image defines several constraints:
// - all glyphs must fit into fixed-size grid;
// - cell size of the grid is estimated by advance width and height of most
//   glyphs in the font;
// - double-width glyphs are skipped;
// - rasterized glyphs that doesn't fit into cell size will be clipped;
func (font *Font) raster(
	dpi float64,
	size float64,
	hinting xfont.Hinting,
) error {
	context := freetype.NewContext()

	context.SetDPI(dpi)
	context.SetFontSize(size)
	context.SetHinting(hinting)
	context.SetFont(font.handle)

	font.estimateMetrics(context.GetScale())

	var (
		cells = int(math.Ceil(math.Sqrt(float64(font.metrics.length))))

		width  = cells * font.metrics.width
		height = cells * font.metrics.height
	)

	font.Image = image.NewRGBA(image.Rect(0, 0, width, height))

	var (
		column = 0
		row    = 0
	)

	for _, segment := range font.handle.Chars() {
		for char := segment.Start; char < segment.End; char++ {
			drawn, err := font.rasterChar(context, row, column, char)
			if err != nil {
				return err
			}

			if drawn {
				font.Glyphs[string(char)] = &Glyph{
					Row:    row,
					Column: column,
					Char:   string(char),
				}

				column++

				if column >= cells {
					row++
					column = 0
				}
			}
		}
	}

	return nil
}

func (font *Font) rasterChar(
	context *freetype.Context,
	row int,
	column int,
	char rune,
) (bool, error) {
	index := font.handle.Index(char)
	width := font.handle.HMetric(context.GetScale(), index).AdvanceWidth.Ceil()

	point := fixed.Point26_6{
		X: fixed.I(font.metrics.width * column),
		Y: fixed.I(font.metrics.height*(row+1) + font.metrics.descender),
	}

	// TODO: Support double width characters, which we skip now.
	if width > font.metrics.width {
		return false, nil
	}

	_, mask, offset, err := context.Glyph(index, point)
	if err != nil {
		return false, karma.Format(
			err,
			")unable to raster glyph %c",
			char,
		)
	}

	// Some glyphs be bigger, than cell in our rendering grid, so we need to
	// clip them into cell size.
	var (
		cell = image.Rect(
			font.metrics.width*column,
			font.metrics.height*row,
			font.metrics.width*(column+1),
			font.metrics.height*(row+1),
		)

		bounds = mask.Bounds().Add(offset)
		box    = cell.Intersect(bounds)
	)

	if box.Empty() {
		return false, nil
	}

	// Glyphs that doesn't fit cell completely must be aligned properly, so
	// clipping will take place either at top or at left side of glyph image.
	pivot := image.Pt(
		cell.Min.X-bounds.Min.X,
		cell.Min.Y-bounds.Min.Y,
	)

	if pivot.X < 0 {
		pivot.X = 0
	}

	if pivot.Y < 0 {
		pivot.Y = 0
	}

	draw.DrawMask(
		font.Image, box,
		image.Black, image.ZP,
		mask, pivot,
		draw.Over,
	)

	return true, nil
}

// estimateMetrics loops over every glyph in font to calculate following
// metrics:
// - total amount of characters defined in font;
// - advance width of most characters in font;
// - advance height of most characters in font;
// - font descender if any.
func (font *Font) estimateMetrics(scale fixed.Int26_6) {
	length := 0

	var (
		widths  = map[int]int{}
		heights = map[int]int{}
	)

	for _, segment := range font.handle.Chars() {
		length += int(segment.End - segment.Start)

		for char := segment.Start; char < segment.End; char++ {
			var (
				index  = font.handle.Index(char)
				width  = font.handle.HMetric(scale, index).AdvanceWidth.Ceil()
				height = font.handle.VMetric(scale, index).AdvanceHeight.Ceil()
			)

			widths[width]++
			heights[height]++
		}
	}

	most := func(measures map[int]int) int {
		var (
			max    = 0
			result = 0
		)

		for measure, count := range measures {
			if count > max {
				max = count
				result = measure
			}
		}

		return result
	}

	font.metrics.length = length
	font.metrics.width = most(widths)
	font.metrics.height = most(heights)

	font.metrics.descender = font.handle.GetDescender().Ceil()
}

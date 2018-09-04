package fonts

import (
	"archive/tar"
	"image"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kovetskiy/toml"
	"github.com/reconquest/karma-go"

	"image/color"
	_ "image/png"
)

type Meta struct {
	Version string
	Width   int
	Height  int
}

type Glyph struct {
	Row    int
	Column int
	Char   string
}

type Font struct {
	Image  *image.RGBA
	Index  []*Glyph
	Glyphs map[string]*Glyph
	Meta   *Meta
}

func Load(path string) (*Font, error) {
	description := karma.Describe("path", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, description.Format(
			err,
			"{font} unable to load font",
		)
	}

	reader := tar.NewReader(file)

	meta, index, glyphs, err := readComponents(reader)
	if err != nil {
		return nil, description.Format(
			err,
			"{font} unable to decode component",
		)
	}

	if meta.Version != "1" {
		return nil, description.Format(
			err,
			"{font} version is not supported: %s",
			meta.Version,
		)
	}

	if len(index) == 0 {
		return nil, description.Format(
			err,
			"{font} empty index in font",
		)
	}

	var (
		columns = glyphs.Bounds().Dx() / meta.Width
		rows    = glyphs.Bounds().Dy() / meta.Height
	)

	if len(index)/columns > rows {
		return nil, description.Reason(
			"{font} index contains more glyphs than image",
		)
	}

	font := &Font{
		Image:  processGlyphsImage(glyphs),
		Glyphs: map[string]*Glyph{},
		Meta:   meta,
	}

	for i, char := range index {
		glyph := &Glyph{
			Row:    i / columns,
			Column: i % columns,
			Char:   char,
		}

		font.Index = append(font.Index, glyph)
		font.Glyphs[char] = glyph
	}

	return font, nil
}

func readComponents(
	reader *tar.Reader,
) (*Meta, []string, image.Image, error) {
	var (
		meta   Meta
		index  []string
		glyphs image.Image
	)

	for {
		entry, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, nil, nil, karma.Format(
				err,
				"malformed tar",
			)
		}

		switch entry.Name {
		case "meta":
			_, err = toml.DecodeReader(reader, &meta)
			if err != nil {
				return nil, nil, nil, karma.Format(
					err,
					"unable to decode font meta",
				)
			}

		case "index":
			body, err := ioutil.ReadAll(reader)
			if err != nil {
				return nil, nil, nil, karma.Format(
					err,
					"unable to decode font index",
				)
			}

			index = strings.Split(string(body), "\n")

		case "glyphs":
			glyphs, _, err = image.Decode(reader)
			if err != nil {
				return nil, nil, nil, karma.Format(
					err,
					"unable to decode font glyphs",
				)
			}
		}
	}

	return &meta, index, glyphs, nil
}

func processGlyphsImage(input image.Image) *image.RGBA {
	glyphs := image.NewRGBA(input.Bounds())

	width, height := input.Bounds().Dx(), input.Bounds().Dy()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, b, g, a := input.At(x, y).RGBA()

			if a > 0 && r == g && g == b && r < 255 {
				glyphs.Set(x, y, color.White)
			}
		}
	}

	return glyphs
}

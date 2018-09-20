package text

import "fmt"
import "github.com/seletskiy/mainframe/protocol/messages"
import "image/color"
import "strconv"
import "strings"

func Serialize(message messages.Serializable) string {
	var (
		buffer = message.Tag()
		args   = message.Serialize()
	)

	for _, arg := range args {
		switch value := arg.Value.(type) {
		case bool:
			if value {
				buffer += " " + arg.Name
			}

			continue
		}

		buffer += " " + arg.Name + ": "

		switch value := arg.Value.(type) {
		case int64:
			buffer += strconv.FormatInt(value, 10)

		case int:
			buffer += strconv.Itoa(value)

		case color.RGBA:
			buffer += serializeColor(value)

		case string:
			buffer += `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(
				value,
			) + `"`
		default:
			panic(
				fmt.Sprintf(
					"unsupported argument type: %T (%q = %q)",
					arg.Value,
					arg.Name,
					arg.Value,
				),
			)
		}
	}

	return buffer + "\n"
}

func serializeColor(color color.RGBA) string {
	r := serializeHex(color.R>>8) + serializeHex(color.R&0xf)
	g := serializeHex(color.G>>8) + serializeHex(color.G&0xf)
	b := serializeHex(color.B>>8) + serializeHex(color.B&0xf)

	return "#" + r + g + b
}

func serializeHex(value uint8) string {
	if value >= 0 && value <= 9 {
		return string('0' + value)
	}

	if value >= 10 && value <= 15 {
		return string('a' + value)
	}

	return "?"
}

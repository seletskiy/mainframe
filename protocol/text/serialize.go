package text

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/seletskiy/mainframe/protocol/messages"
)

func Serialize(message messages.Serializable) []byte {
	var (
		buffer = message.Tag()
		args   = message.Serialize()
	)

	for _, arg := range args {
		if value, ok := arg.Value.(bool); ok {
			if value {
				buffer += " " + arg.Name
			}

			continue
		}

		if value := serializeValue(arg.Value); value != "" {
			buffer += " " + arg.Name + ": " + value
		}
	}

	return []byte(buffer + "\n")
}

func serializeValue(value interface{}) string {
	switch value := value.(type) {
	case *int64:
		if value == nil {
			return ""
		}

		return serializeValue(*value)

	case *int:
		if value == nil {
			return ""
		}

		return serializeValue(*value)

	case *color.RGBA:
		if value == nil {
			return ""
		}

		return serializeValue(*value)

	case *string:
		if value == nil {
			return ""
		}

		return serializeValue(*value)

	case int64:
		return strconv.FormatInt(value, 10)

	case int:
		return strconv.Itoa(value)

	case color.RGBA:
		return serializeColor(value)

	case string:
		return strconv.Quote(value)

	default:
		panic(
			fmt.Sprintf(
				"unsupported argument type: %[1]T (%[1]v)",
				value,
			),
		)
	}
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

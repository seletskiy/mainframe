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
		if value, flag := serializeValue(arg.Value); value != "" {
			if flag {
				buffer += " " + arg.Name
			} else {
				buffer += " " + arg.Name + ": " + value
			}
		}
	}

	return []byte(buffer + "\n")
}

func serializeValue(value interface{}) (string, bool) {
	switch value := value.(type) {
	case *bool:
		if value == nil {
			return "", true
		}

		return serializeValue(*value)

	case *int64:
		if value == nil {
			return "", false
		}

		return serializeValue(*value)

	case *int:
		if value == nil {
			return "", false
		}

		return serializeValue(*value)

	case *color.RGBA:
		if value == nil {
			return "", false
		}

		return serializeValue(*value)

	case *string:
		if value == nil {
			return "", false
		}

		return serializeValue(*value)

	case bool:
		if value {
			return "true", true
		} else {
			return "", true
		}

	case int64:
		return strconv.FormatInt(value, 10), false

	case int:
		return strconv.Itoa(value), false

	case color.RGBA:
		return serializeColor(value), false

	case string:
		return strconv.Quote(value), false

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

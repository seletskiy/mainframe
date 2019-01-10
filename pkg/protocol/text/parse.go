package text

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/seletskiy/mainframe/pkg/protocol/messages"
)

type ParseFunc func(map[string]interface{}) (messages.Tagged, error)

func Parse(data string) (messages.Tagged, error) {
	var (
		tag  string
		args = map[string]interface{}{}
	)

	for _, token := range tokenize(data) {
		switch {
		case token["tag"] != "":
			tag = token["tag"]

		case token["key"] != "":
			var value interface{}

			switch {
			case token["int"] != "":
				value, _ = strconv.Atoi(token["int"])

			case token["color"] != "":
				value = parseColor(token["color"])

			case token["string"] != "":
				var err error
				value, err = strconv.Unquote(token["string"])
				if err != nil {
					return nil, fmt.Errorf(
						`unable to unquote string: %q`,
						token["string"],
					)
				}

			default:
				value = true
			}

			args[token["key"]] = value

		case token["garbage"] != "":
			return nil, fmt.Errorf(
				`unexpected token: %q`,
				token["garbage"],
			)
		}
	}

	parsers := map[string]ParseFunc{
		"ok":        parseOKMessage,
		"error":     parseErrorMessage,
		"put":       parsePutMessage,
		"subscribe": parseSubscribeMessage,
		"open":      parseOpenMessage,
		"reshape":   parseReshapeMessage,
		"clear":     parseClearMessage,
		"get":       parseGetMessage,
	}

	if parser, ok := parsers[tag]; ok {
		return parser(args)
	} else {
		return nil, fmt.Errorf(`unknown message tag: %s`, tag)
	}
}

func tokenize(data string) []map[string]string {
	var (
		names   = reMessage.SubexpNames()
		matches = reMessage.FindAllStringSubmatch(data, -1)
	)

	tokens := []map[string]string{}

	for _, submatches := range matches {
		token := map[string]string{}

		for i, name := range names[1:] {
			token[name] = submatches[i+1]
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func parseColor(data string) color.RGBA {
	var step int

	switch len(data) {
	case 3:
		step = 1
	case 6:
		step = 2
	}

	r, _ := strconv.ParseInt(data[0*step:1*step], 16, 0)
	g, _ := strconv.ParseInt(data[1*step:2*step], 16, 0)
	b, _ := strconv.ParseInt(data[2*step:3*step], 16, 0)

	if step == 1 {
		r = r + r<<4
		g = g + g<<4
		b = b + b<<4
	}

	return color.RGBA{
		uint8(r),
		uint8(g),
		uint8(b),
		0xff,
	}
}

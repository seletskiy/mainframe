package text

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/seletskiy/mainframe/protocol/messages"
)

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
				value = strings.NewReplacer(`\\`, `\`, `\"`, `"`).Replace(
					token["string"],
				)

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

	var (
		message messages.Tagged
		err     error
	)

	switch tag {
	case "put":
		message, err = parsePutMessage(args)
		if err != nil {
			return nil, err
		}
	}

	return message, nil
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

	return color.RGBA{
		uint8(r),
		uint8(g),
		uint8(b),
		0xff,
	}
}

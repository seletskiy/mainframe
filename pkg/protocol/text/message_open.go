package text

import (
	"fmt"

	"github.com/seletskiy/mainframe/pkg/protocol/messages"
)

func parseOpenMessage(
	args map[string]interface{},
) (messages.Tagged, error) {
	message := &messages.Open{}

	size, err := parseSize(args)
	if err != nil {
		return nil, err
	}

	message.Size = size

	switch {
	case args["x"] != nil && args["y"] == nil:
		fallthrough
	case args["x"] == nil && args["y"] != nil:
		return nil, fmt.Errorf("x and y should be specified together")
	}

	err = NewSpec().
		Skip("width").
		Skip("height").
		Skip("rows").
		Skip("columns").
		Int("x", &message.X).
		Int("y", &message.Y).
		String("title", &message.Title).
		Bool("raw", &message.Raw).
		Bool("hidden", &message.Hidden).
		Bool("fixed", &message.Fixed).
		Bool("bare", &message.Bare).
		Bool("floating", &message.Floating).
		Bind(args)
	if err != nil {
		return nil, err
	}

	return message, nil
}

package text

import (
	"fmt"

	"github.com/seletskiy/mainframe/protocol/messages"
)

func parseOpenMessage(
	args map[string]interface{},
) (messages.Tagged, error) {
	message := &messages.Open{}

	switch {
	case args["width"] != nil && args["height"] == nil:
		fallthrough
	case args["width"] == nil && args["height"] != nil:
		return nil, fmt.Errorf("width and height should be specified together")
	}

	switch {
	case args["x"] != nil && args["y"] == nil:
		fallthrough
	case args["x"] == nil && args["y"] != nil:
		return nil, fmt.Errorf("x and y should be specified together")
	}

	err := NewSpec().
		Int("width", &message.Width).
		Int("height", &message.Height).
		Int("x", &message.X).
		Int("y", &message.Y).
		String("title", &message.Title).
		Bool("raw", &message.Raw).
		Bool("hidden", &message.Hidden).
		Bool("fixed", &message.Fixed).
		Bool("decorated", &message.Decorated).
		Bool("floating", &message.Floating).
		Bind(args)
	if err != nil {
		return nil, err
	}

	if message.Width != nil {
		if *message.Width <= 0 {
			return nil, fmt.Errorf("width should be greater than zero")
		}
	}

	if message.Height != nil {
		if *message.Height <= 0 {
			return nil, fmt.Errorf("height should be greater than zero")
		}
	}

	return message, nil
}

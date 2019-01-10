package text

import (
	"fmt"

	"github.com/seletskiy/mainframe/pkg/protocol/messages"
)

func parseReshapeMessage(
	args map[string]interface{},
) (messages.Tagged, error) {
	message := &messages.Reshape{}

	size, err := parseSize(args)
	if err != nil {
		return nil, err
	}

	message.Size = size

	switch {
	case args["x"] != nil && args["y"] == nil:
		fallthrough
	case args["x"] == nil && args["y"] != nil:
		return nil, fmt.Errorf("x & y should be specified together")
	}

	switch {
	case args["x"] != nil && args["y"] != nil:
	case args["width"] != nil && args["height"] != nil:
	default:
		return nil, fmt.Errorf("x & y or width & height should be specified")
	}

	err = NewSpec().
		Skip("rows").
		Skip("columns").
		Int("x", &message.X).
		Int("y", &message.Y).
		Bind(args)
	if err != nil {
		return nil, err
	}

	return message, nil
}

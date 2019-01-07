package text

import (
	"fmt"

	"github.com/seletskiy/mainframe/protocol/messages"
)

func parseClearMessage(
	args map[string]interface{},
) (messages.Tagged, error) {
	message := &messages.Clear{}

	switch {
	case args["x"] != nil && args["y"] == nil:
		fallthrough
	case args["x"] == nil && args["y"] != nil:
		return nil, fmt.Errorf("x & y should be specified together")
	}

	err := NewSpec().
		Int("x", &message.X).
		Int("y", &message.Y).
		Int("columns", &message.Columns).
		Int("rows", &message.Rows).
		Bind(args)
	if err != nil {
		return nil, err
	}

	return message, nil
}

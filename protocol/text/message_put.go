package text

import (
	"github.com/seletskiy/mainframe/protocol/messages"
)

func parsePutMessage(
	args map[string]interface{},
) (messages.Tagged, error) {
	switch {
	case args["fg"] != nil:
	case args["bg"] != nil:
	case args["text"] != nil:
	default:
		return nil, ErrMissingGroup{"fg", "bg", "text"}
	}

	message := &messages.Put{}

	err := NewSpec().
		Require("x").
		Require("y").
		Int("x", &message.X).
		Int("y", &message.Y).
		Int("columns", &message.Columns).
		Int("rows", &message.Rows).
		Color("fg", &message.Foreground).
		Color("bg", &message.Background).
		String("text", &message.Text).
		Int("tick", &message.Tick).
		Bool("exclusive", &message.Exclusive).
		Bind(args)
	if err != nil {
		return nil, err
	}

	return message, nil
}

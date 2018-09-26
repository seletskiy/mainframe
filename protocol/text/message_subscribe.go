package text

import (
	"github.com/seletskiy/mainframe/protocol/messages"
)

func parseSubscribeMessage(
	args map[string]interface{},
) (messages.Tagged, error) {
	message := &messages.Subscribe{}

	err := NewSpec().
		Bool("resize", &message.Resize).
		Bool("keyboard", &message.Keyboard).
		Bool("input", &message.Input).
		Bind(args)
	if err != nil {
		return nil, err
	}

	return message, nil
}

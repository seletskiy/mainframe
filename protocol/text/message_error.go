package text

import (
	"github.com/seletskiy/mainframe/protocol/messages"
)

func parseErrorMessage(
	args map[string]interface{},
) (messages.Tagged, error) {
	message := &messages.Error{}

	err := NewSpec().
		Require("message").
		String("message", &message.Message).
		Bind(args)
	if err != nil {
		return nil, err
	}

	return message, nil
}

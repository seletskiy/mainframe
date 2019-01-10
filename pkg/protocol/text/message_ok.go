package text

import (
	"github.com/seletskiy/mainframe/pkg/protocol/messages"
)

func parseOKMessage(
	args map[string]interface{},
) (messages.Tagged, error) {
	message := &messages.OK{}

	for name, value := range args {
		message.Set(name, value)
	}

	return message, nil
}

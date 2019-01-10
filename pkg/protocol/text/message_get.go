package text

import (
	"fmt"

	"github.com/seletskiy/mainframe/pkg/protocol/messages"
)

func parseGetMessage(
	args map[string]interface{},
) (messages.Tagged, error) {
	message := &messages.Get{}

	err := NewSpec().
		SkipUnknown().
		Bool("font", &message.Font.Set).
		Bind(args)
	if err != nil {
		return nil, err
	}

	switch {
	case message.Font.Set:
		//
	default:
		return nil, fmt.Errorf(
			"no attribute specified for get",
		)
	}

	return message, nil
}

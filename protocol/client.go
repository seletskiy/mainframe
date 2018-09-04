package protocol

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/seletskiy/mainframe/engine"
)

type Client struct {
	connection net.Conn

	context *engine.Context
}

func (client *Client) serve() {
	scanner := bufio.NewScanner(client.connection)

	for scanner.Scan() {
		components := strings.SplitN(scanner.Text(), " ", 2)
		if len(components) != 2 {
			continue
		}

		var (
			command   = strings.ToLower(components[0])
			arguments = components[1]
		)

		switch command {
		case "text":
			client.handleText(arguments)
		}
	}

	client.context.Close()
}

func (client *Client) handleText(arguments string) error {
	components := strings.SplitN(arguments, " ", 3)
	if len(components) != 3 {
		return fmt.Errorf("text command expect to have 3 arguments")
	}

	x, err := strconv.Atoi(components[0])
	if err != nil {
		return fmt.Errorf("unable to parse X coordinate: %s", err)
	}

	y, err := strconv.Atoi(components[1])
	if err != nil {
		return fmt.Errorf("unable to parse Y coordinate: %s", err)
	}

	text := components[2]

	client.context.Screen.SetText(x, y, text)

	return nil
}

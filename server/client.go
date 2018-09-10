package server

import (
	"bufio"
	"net"

	"github.com/seletskiy/mainframe/engine"
	"github.com/seletskiy/mainframe/protocol/messages"
	"github.com/seletskiy/mainframe/protocol/text"
)

type Client struct {
	connection net.Conn

	context *engine.Context
}

func (client *Client) serve() {
	scanner := bufio.NewScanner(client.connection)

	for scanner.Scan() {
		line := scanner.Text()

		tag, message, err := text.Parse(line)
		if err != nil {
			panic(err)
		}

		switch tag {
		case "put":
			message := message.(*messages.Put)
			client.context.Screen.SetText(message.X, message.Y, *message.Text)
		}
	}

	client.context.Close()
}

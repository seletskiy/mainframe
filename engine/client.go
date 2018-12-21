package engine

import (
	"bufio"
	"net"

	"github.com/reconquest/karma-go"
	"github.com/seletskiy/mainframe/protocol/messages"
	"github.com/seletskiy/mainframe/protocol/text"
)

type Client struct {
	Connection net.Conn
	Context    *Context
	Engine     *Engine
}

func (client *Client) Serve() {
	scanner := bufio.NewScanner(client.Connection)

	for scanner.Scan() {
		line := scanner.Text()

		message, err := text.Parse(line)
		if err != nil {
			err = client.Error(err)
			if err != nil {
				// FIXME
				panic(err)
			}
		}

		switch message := message.(type) {
		case *messages.Put:
			err = client.handlePut(message)

		case *messages.Subscribe:
			err = client.handleSubscribe(message)

		case *messages.Open:
			err = client.handleOpen(message)

		case *messages.Get:
			err = client.handleGet(message)
		}

		if err != nil {
			err = client.Error(err)
			if err != nil {
				// FIXME
				panic(err)
			}
		}
	}

	if client.Context != nil {
		client.Context.Close()
	}
}

func (client *Client) Send(message messages.Serializable) error {
	_, err := client.Connection.Write([]byte(text.Serialize(message)))
	if err != nil {
		return karma.
			Describe("message", message).
			Format(
				err,
				"unable to send message to client",
			)
	}

	return nil
}

func (client *Client) Error(err error) error {
	return client.Send(&messages.Error{
		Message: err.Error(),
	})
}

func (client *Client) handlePut(message *messages.Put) error {
	var reply messages.OK

	if !client.Context.screen.Put(message) {
		reply.Set("offscreen", true)
	}

	return client.Send(&reply)
}

func (client *Client) handleSubscribe(message *messages.Subscribe) error {
	var reply messages.OK

	if message.Resize {
		client.Context.Subscribe(client, SubscriptionResize)
	}

	if message.Keyboard {
		client.Context.Subscribe(client, SubscriptionKeyboard)
	}

	if message.Input {
		client.Context.Subscribe(client, SubscriptionInput)
	}

	return client.Send(&reply)
}

func (client *Client) handleOpen(message *messages.Open) error {
	var err error

	client.Context, err = client.Engine.CreateWindow(message)
	if err != nil {
		return err
	}

	var reply messages.OK

	// TODO: reply with id of newly created window
	//reply.Set("id", 123)

	return client.Send(&reply)
}

func (client *Client) handleGet(message *messages.Get) error {
	var reply messages.OK

	switch {
	case message.Font.Set:
		font := client.Engine.GetFont()

		reply.Set("width", font.Meta.Width)
		reply.Set("height", font.Meta.Height)
	}

	return client.Send(&reply)
}

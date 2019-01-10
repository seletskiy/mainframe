package engine

import (
	"bufio"
	"net"

	"github.com/reconquest/karma-go"
	"github.com/seletskiy/mainframe/pkg/protocol/messages"
	"github.com/seletskiy/mainframe/pkg/protocol/text"
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

		case *messages.Reshape:
			err = client.handleReshape(message)

		case *messages.Clear:
			err = client.handleClear(message)
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
			Describe("message", message.Serialize()).
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

	if !client.Context.Screen.Put(message) {
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

		reply.Set("width", font.GetWidth())
		reply.Set("height", font.GetHeight())
	}

	return client.Send(&reply)
}

func (client *Client) handleReshape(message *messages.Reshape) error {
	var reply messages.OK

	font := client.Engine.GetFont()

	var (
		width  int
		height int
	)

	if message.Width != nil {
		width = *message.Width
	}

	if message.Height != nil {
		height = *message.Height
	}

	if message.Columns != nil {
		width = *message.Columns * font.GetWidth()
	}

	if message.Rows != nil {
		height = *message.Rows * font.GetHeight()
	}

	if message.X != nil && message.Y != nil {
		client.Context.Window.SetPos(*message.X, *message.Y)
	}

	if width > 0 && height > 0 {
		client.Context.Window.SetSize(width, height)
		client.Context.Resize(width, height)
	}

	return client.Send(&reply)
}

func (client *Client) handleClear(message *messages.Clear) error {
	var reply messages.OK

	var (
		x int
		y int

		rows    int
		columns int
	)

	if message.X != nil {
		x = *message.X
	}

	if message.Y != nil {
		y = *message.Y
	}

	if message.Rows != nil {
		rows = *message.Rows
	} else {
		rows = client.Context.Screen.rows
	}

	if message.Columns != nil {
		columns = *message.Columns
	} else {
		_, columns = client.Context.Screen.GetGrid()
	}

	client.Context.Screen.Clear(x, y, rows, columns)

	return client.Send(&reply)
}
